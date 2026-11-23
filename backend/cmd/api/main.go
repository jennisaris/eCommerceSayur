package main

import (
	"ecommerce-backend/internal/config"
	deliveryHttp "ecommerce-backend/internal/delivery/http"
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/pkg/email"
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/midtrans"
	"ecommerce-backend/internal/pkg/seclog"
	postgresrepo "ecommerce-backend/internal/repository/postgres"
	redisrepo "ecommerce-backend/internal/repository/redis"
	"ecommerce-backend/internal/usecase"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("🚀 Starting eCommerce High-Performance Backend Service...")

	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	seclog.ConfigureAlertWebhook(os.Getenv("SECURITY_ALERT_WEBHOOK_URL"))

	// 2. Initialize PostgreSQL Connection
	db, err := postgresrepo.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("✅ PostgreSQL connected successfully")

	// 3. Run Migrations (Check both /app/migrations and ./migrations)
	migrationsDir := "/app/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join(".", "migrations")
	}
	if err := postgresrepo.RunMigrations(db, migrationsDir, cfg.RunSeeder); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// 4. Initialize Redis Connection (with graceful degradation if offline)
	rdb, err := redisrepo.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("⚠️ Redis connection warning (caching will be disabled): %v", err)
	} else {
		log.Println("✅ Redis cache connected successfully")
	}

	// 5. Initialize Infrastructure Repositories
	userRepo := postgresrepo.NewUserRepository(db)
	categoryRepo := postgresrepo.NewCategoryRepository(db)
	productRepo := postgresrepo.NewProductRepository(db)
	cartRepo := postgresrepo.NewCartRepository(db)
	orderRepo := postgresrepo.NewOrderRepository(db)
	paymentRepo := postgresrepo.NewPaymentRepository(db)
	notifRepo := postgresrepo.NewNotificationRepository(db)

	var catalogCache *redisrepo.CatalogCacheRepository
	var cartCache *redisrepo.CartCacheRepository
	if rdb != nil {
		catalogCache = redisrepo.NewCatalogCacheRepository(rdb)
		cartCache = redisrepo.NewCartCacheRepository(rdb)
	}

	// 6. Initialize Utility Services
	jwtSvc := jwt.NewJWTServiceWithDuration(cfg.JWT.Secret, time.Duration(cfg.JWT.AccessExpirationMinutes)*time.Minute)
	midtransSvc := midtrans.NewClient(cfg.Midtrans.ServerKey, cfg.Midtrans.ClientKey, cfg.Midtrans.IsProduction)

	// 7. Initialize Use Cases (Clean Architecture)
	smtpSender := email.NewSMTPSender(email.SMTPConfig{
		Host: cfg.Email.SMTPHost, Port: cfg.Email.SMTPPort, Username: cfg.Email.SMTPUsername,
		Password: cfg.Email.SMTPPassword, From: cfg.Email.From,
	})
	authUseCase := usecase.NewAuthUseCaseWithOptions(userRepo, cartRepo, jwtSvc, usecase.AuthOptions{
		RefreshTokenTTL:          time.Duration(cfg.JWT.RefreshTokenDays) * 24 * time.Hour,
		RequireEmailVerification: cfg.RequireEmailVerification,
		VerificationEmailSender:  smtpSender,
		PasswordResetEmailSender: smtpSender,
		VerificationBaseURL:      cfg.Client.BaseURL,
		IsProduction:             cfg.AppEnv == "production",
	})
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, catalogCache)
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, catalogCache)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo, cartCache)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo, userRepo, notifRepo, cartCache)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, orderRepo, userRepo, notifRepo, midtransSvc)

	// 8. Initialize HTTP Handlers
	authHandler := handler.NewAuthHandler(
		authUseCase,
		cfg.JWT.CookieSecure,
		time.Duration(cfg.JWT.AccessExpirationMinutes)*time.Minute,
		time.Duration(cfg.JWT.RefreshTokenDays)*24*time.Hour,
	)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	adminOrderHandler := handler.NewAdminOrderHandler(orderUseCase)

	// 9. Initialize Fiber Web Application
	app := fiber.New(fiber.Config{
		AppName:                 "BookVerse - Enterprise BookStore Platform API",
		ServerHeader:            "BookVerse-Go-CleanArch",
		BodyLimit:               1 * 1024 * 1024, // 1MB
		ReadBufferSize:          8192,
		WriteBufferSize:         8192,
		DisableStartupMessage:   false,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		ProxyHeader:             "X-Real-IP",
	})

	// 10. Setup Routes
	deliveryHttp.SetupRoutes(app, &deliveryHttp.RouterConfig{
		AuthHandler:       authHandler,
		CategoryHandler:   categoryHandler,
		ProductHandler:    productHandler,
		CartHandler:       cartHandler,
		OrderHandler:      orderHandler,
		PaymentHandler:    paymentHandler,
		AdminOrderHandler: adminOrderHandler,
		UserRepository:    userRepo,
		JWTService:        jwtSvc,
		AllowedOrigins:    cfg.CORSAllowedOrigins,
		CookieSecure:      cfg.JWT.CookieSecure,
		IsProduction:      cfg.AppEnv == "production",
	})

	// 11. Graceful Shutdown Handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Gracefully shutting down eCommerce API server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server listening and serving on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}

package main

import (
	"ecommerce-backend/internal/config"
	deliveryHttp "ecommerce-backend/internal/delivery/http"
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/midtrans"
	postgresrepo "ecommerce-backend/internal/repository/postgres"
	redisrepo "ecommerce-backend/internal/repository/redis"
	"ecommerce-backend/internal/usecase"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("🚀 Starting eCommerce High-Performance Backend Service...")

	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

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
	if err := postgresrepo.RunMigrations(db, migrationsDir); err != nil {
		log.Printf("⚠️ Migration warning: %v", err)
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
	jwtSvc := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	midtransSvc := midtrans.NewClient(cfg.Midtrans.ServerKey, cfg.Midtrans.ClientKey, cfg.Midtrans.IsProduction)

	// 7. Initialize Use Cases (Clean Architecture)
	authUseCase := usecase.NewAuthUseCase(userRepo, cartRepo, jwtSvc)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, catalogCache)
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, catalogCache)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo, cartCache)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo, userRepo, notifRepo, cartCache)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, orderRepo, userRepo, notifRepo, midtransSvc)

	// 8. Initialize HTTP Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	adminOrderHandler := handler.NewAdminOrderHandler(orderUseCase)

	// 9. Initialize Fiber Web Application
	app := fiber.New(fiber.Config{
		AppName:               "Enterprise eCommerce Platform API",
		ServerHeader:          "Fiber-Go-CleanArch",
		BodyLimit:             10 * 1024 * 1024, // 10MB
		ReadBufferSize:        8192,
		WriteBufferSize:       8192,
		DisableStartupMessage: false,
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
		JWTService:        jwtSvc,
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

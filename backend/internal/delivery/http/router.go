package http

import (
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/delivery/http/middleware"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouterConfig struct {
	AuthHandler       *handler.AuthHandler
	CategoryHandler   *handler.CategoryHandler
	ProductHandler    *handler.ProductHandler
	CartHandler       *handler.CartHandler
	OrderHandler      *handler.OrderHandler
	PaymentHandler    *handler.PaymentHandler
	AdminOrderHandler *handler.AdminOrderHandler
	UserRepository    domain.UserRepository
	JWTService        *jwt.JWTService
	AllowedOrigins    string
	CookieSecure      bool
	IsProduction      bool
}

func newSecurityHeadersMiddleware(isProduction bool) fiber.Handler {
	cfg := helmet.Config{
		XSSProtection:         "0",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
	}
	if isProduction {
		cfg.HSTSMaxAge = 31536000
		cfg.HSTSExcludeSubdomains = false
	}
	helmetHandler := helmet.New(cfg)
	return func(c *fiber.Ctx) error {
		if isProduction {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		return helmetHandler(c)
	}
}

func SetupRoutes(app *fiber.App, cfg *RouterConfig) {
	// Global Middlewares
	app.Use(recover.New())
	app.Use(newSecurityHeadersMiddleware(cfg.IsProduction))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	allowedOrigins := cfg.AllowedOrigins
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://127.0.0.1:3000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		AllowCredentials: true,
	}))
	app.Use(middleware.CSRFProtection(cfg.CookieSecure, allowedOrigins))

	// Security Rate Limiters
	// 1. Auth Limiter: 10 attempts per minute per IP (login, register, refresh, verify-email)
	authLimiter := newIPRateLimiter(10, 1*time.Minute, "Terlalu banyak percobaan autentikasi. Silakan tunggu 1 menit.")

	// 2. Password Reset Limiter: 5 attempts per 15 minutes per IP (prevents SMTP abuse & token brute-force)
	passwordResetLimiter := newIPRateLimiter(5, 15*time.Minute, "Terlalu banyak permintaan reset kata sandi. Silakan coba lagi nanti.")

	// 3. Order Checkout & Estimate Limiter: 30 attempts per minute per IP
	orderLimiter := newIPRateLimiter(30, 1*time.Minute, "Terlalu banyak permintaan pemesanan. Silakan tunggu sesaat.")

	// 4. Payment Gateway (Midtrans Snap) Limiter: 15 attempts per minute per user/IP
	paymentLimiter := newUserOrIPRateLimiter(15, 1*time.Minute, "Terlalu banyak permintaan pembayaran. Silakan tunggu 1 menit.")

	// Root Welcome & Health Check Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"service":  "BookVerse - Enterprise BookStore Platform API",
			"status":   "online",
			"version":  "v1.0.0",
			"frontend": "http://127.0.0.1:3000",
			"health":   "/health",
			"api_v1":   "/api/v1",
			"message":  "Backend BookVerse API aktif. Untuk antarmuka web toko buku pelanggan & admin, silakan buka http://127.0.0.1:3000",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "bookverse-api",
		})
	})

	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"service": "BookVerse REST API v1",
			"status":  "running",
			"endpoints": fiber.Map{
				"health":     "/health",
				"auth":       "/api/v1/auth",
				"categories": "/api/v1/categories",
				"products":   "/api/v1/products",
				"cart":       "/api/v1/cart",
				"orders":     "/api/v1/orders",
				"payments":   "/api/v1/payments",
			},
		})
	})

	// ==========================================
	// PUBLIC ROUTES
	// ==========================================
	// Auth (Protected with Rate Limiter)
	auth := api.Group("/auth")
	auth.Get("/csrf", func(c *fiber.Ctx) error {
		token, err := middleware.CSRFToken(c)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Unable to initialize security token")
		}
		return response.Success(c, "CSRF token ready", fiber.Map{"token": token})
	})
	auth.Post("/verify-email", authLimiter, cfg.AuthHandler.VerifyEmail)
	auth.Post("/register", authLimiter, cfg.AuthHandler.Register)
	auth.Post("/login", authLimiter, cfg.AuthHandler.Login)
	auth.Post("/refresh", authLimiter, cfg.AuthHandler.Refresh)
	auth.Post("/forgot-password", passwordResetLimiter, cfg.AuthHandler.ForgotPassword)
	auth.Post("/reset-password", passwordResetLimiter, cfg.AuthHandler.ResetPassword)
	auth.Post("/logout", cfg.AuthHandler.Logout)

	// Catalog (Public, Cache-Accelerated)
	api.Get("/categories", cfg.CategoryHandler.GetAll)
	api.Get("/categories/:slug", cfg.CategoryHandler.GetBySlug)
	api.Get("/products", cfg.ProductHandler.GetAll)
	api.Get("/products/:slug", cfg.ProductHandler.GetBySlug)

	// Midtrans Webhook Callback (Public webhook from Midtrans servers)
	api.Post("/payments/midtrans-callback", cfg.PaymentHandler.MidtransCallback)

	// ==========================================
	// AUTHENTICATED USER ROUTES
	// ==========================================
	userAuth := middleware.AuthMiddleware(cfg.JWTService)

	// Profile
	api.Get("/auth/me", userAuth, cfg.AuthHandler.GetMe)
	api.Put("/auth/me", userAuth, authLimiter, cfg.AuthHandler.UpdateProfile)

	// Cart (cross-device sync with Redis cache)
	cart := api.Group("/cart", userAuth)
	cart.Get("/", cfg.CartHandler.GetCart)
	cart.Post("/items", cfg.CartHandler.AddToCart)
	cart.Put("/items/:product_id", cfg.CartHandler.UpdateItem)
	cart.Delete("/items/:product_id", cfg.CartHandler.RemoveItem)
	cart.Delete("/", cfg.CartHandler.ClearCart)

	// Checkout & Orders (100% Delivery Model)
	orders := api.Group("/orders", userAuth)
	orders.Post("/estimate", orderLimiter, cfg.OrderHandler.EstimateCheckout)
	orders.Post("/", orderLimiter, cfg.OrderHandler.CreateOrder)
	orders.Post("/checkout", orderLimiter, cfg.OrderHandler.CreateOrder)
	orders.Get("/", cfg.OrderHandler.GetUserOrders)
	orders.Get("/:code", cfg.OrderHandler.GetOrderByCode)
	orders.Get("/*", cfg.OrderHandler.GetOrderByCode)

	// Payments (rate-limited external Midtrans Snap & Status Sync calls)
	api.Post("/payments/sync", userAuth, paymentLimiter, cfg.PaymentHandler.SyncPaymentStatus)
	api.Post("/payments/:order_code/sync", userAuth, paymentLimiter, cfg.PaymentHandler.SyncPaymentStatus)
	api.Post("/payments/snap", userAuth, paymentLimiter, cfg.PaymentHandler.CreateSnapToken)
	api.Post("/payments/:order_code/snap", userAuth, paymentLimiter, cfg.PaymentHandler.CreateSnapToken)
	api.Post("/payments/*", userAuth, paymentLimiter, cfg.PaymentHandler.CreateSnapToken)

	// ==========================================
	// SUPER ADMIN ROUTES (RBAC Protected)
	// ==========================================
	admin := api.Group("/admin", userAuth, middleware.RequireRole("Super Admin", cfg.UserRepository))

	// Category Management
	admin.Get("/categories", cfg.CategoryHandler.AdminGetAll)
	admin.Post("/categories", cfg.CategoryHandler.Create)
	admin.Put("/categories/:id", cfg.CategoryHandler.Update)
	admin.Delete("/categories/:id", cfg.CategoryHandler.Delete)

	// Product Management
	admin.Post("/products", cfg.ProductHandler.Create)
	admin.Get("/products", cfg.ProductHandler.AdminGetAll)
	admin.Put("/products/:id", cfg.ProductHandler.Update)
	admin.Delete("/products/:id", cfg.ProductHandler.Delete)

	// Order & Shipment Management
	admin.Get("/orders", cfg.AdminOrderHandler.GetAllOrders)
	admin.Put("/orders/:id/status", cfg.AdminOrderHandler.UpdateOrderStatus)
}

func newIPRateLimiter(max int, expiration time.Duration, message string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			seclog.RecordRateLimit(c.IP(), 0, c.Method(), c.Path())
			return response.Error(c, fiber.StatusTooManyRequests, message)
		},
	})
}

func newUserOrIPRateLimiter(max int, expiration time.Duration, message string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID, ok := c.Locals("user_id").(int64); ok && userID > 0 {
				return fmt.Sprintf("user:%d", userID)
			}
			return "ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			var userID int64
			if id, ok := c.Locals("user_id").(int64); ok {
				userID = id
			}
			seclog.RecordRateLimit(c.IP(), userID, c.Method(), c.Path())
			return response.Error(c, fiber.StatusTooManyRequests, message)
		},
	})
}

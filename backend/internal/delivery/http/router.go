package http

import (
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/delivery/http/middleware"
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	JWTService        *jwt.JWTService
	AllowedOrigins    string
}

func SetupRoutes(app *fiber.App, cfg *RouterConfig) {
	// Global Middlewares
	app.Use(recover.New())
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
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
	}))

	// Security Rate Limiters
	// 1. Auth Limiter: 10 attempts per minute per IP to protect from brute-force & credential stuffing
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, fiber.StatusTooManyRequests, "Terlalu banyak percobaan autentikasi. Silakan tunggu 1 menit.")
		},
	})

	// 2. Order Checkout Limiter: 30 attempts per minute per IP
	orderLimiter := limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, fiber.StatusTooManyRequests, "Terlalu banyak permintaan pemesanan. Silakan tunggu sesaat.")
		},
	})

	// Root Welcome & Health Check Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"service":    "Enterprise eCommerce Platform API (Sayur FreshCart)",
			"status":     "online",
			"version":    "v1.0.0",
			"frontend":   "http://127.0.0.1:3000",
			"health":     "/health",
			"api_v1":     "/api/v1",
			"message":    "Backend API aktif. Untuk antarmuka web toko pelanggan & admin, silakan buka http://127.0.0.1:3000",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "ecommerce-api",
		})
	})

	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"service": "Sayur FreshCart REST API v1",
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
	auth.Post("/register", authLimiter, cfg.AuthHandler.Register)
	auth.Post("/login", authLimiter, cfg.AuthHandler.Login)

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

	// Cart (cross-device sync with Redis cache)
	cart := api.Group("/cart", userAuth)
	cart.Get("/", cfg.CartHandler.GetCart)
	cart.Post("/items", cfg.CartHandler.AddToCart)
	cart.Put("/items/:product_id", cfg.CartHandler.UpdateItem)
	cart.Delete("/items/:product_id", cfg.CartHandler.RemoveItem)
	cart.Delete("/", cfg.CartHandler.ClearCart)

	// Checkout & Orders (100% Delivery Model)
	orders := api.Group("/orders", userAuth)
	orders.Post("/estimate", cfg.OrderHandler.EstimateCheckout)
	orders.Post("/", orderLimiter, cfg.OrderHandler.CreateOrder)
	orders.Post("/checkout", orderLimiter, cfg.OrderHandler.CreateOrder)
	orders.Get("/", cfg.OrderHandler.GetUserOrders)
	orders.Get("/:code", cfg.OrderHandler.GetOrderByCode)
	orders.Get("/*", cfg.OrderHandler.GetOrderByCode)

	// Payments
	api.Post("/payments/snap", userAuth, cfg.PaymentHandler.CreateSnapToken)
	api.Post("/payments/:order_code/snap", userAuth, cfg.PaymentHandler.CreateSnapToken)
	api.Post("/payments/*", userAuth, cfg.PaymentHandler.CreateSnapToken)

	// ==========================================
	// SUPER ADMIN ROUTES (RBAC Protected)
	// ==========================================
	admin := api.Group("/admin", userAuth, middleware.RequireRole("Super Admin"))

	// Category Management
	admin.Post("/categories", cfg.CategoryHandler.Create)
	admin.Put("/categories/:id", cfg.CategoryHandler.Update)
	admin.Delete("/categories/:id", cfg.CategoryHandler.Delete)

	// Product Management
	admin.Post("/products", cfg.ProductHandler.Create)
	admin.Put("/products/:id", cfg.ProductHandler.Update)
	admin.Delete("/products/:id", cfg.ProductHandler.Delete)

	// Order & Shipment Management
	admin.Get("/orders", cfg.AdminOrderHandler.GetAllOrders)
	admin.Put("/orders/:id/status", cfg.AdminOrderHandler.UpdateOrderStatus)
}

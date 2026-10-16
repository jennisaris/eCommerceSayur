package http

import (
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/delivery/http/middleware"
	"ecommerce-backend/internal/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
}

func SetupRoutes(app *fiber.App, cfg *RouterConfig) {
	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
	}))

	// Liveness & Readiness Probe
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "ecommerce-api",
		})
	})

	api := app.Group("/api/v1")

	// ==========================================
	// PUBLIC ROUTES
	// ==========================================
	// Auth
	auth := api.Group("/auth")
	auth.Post("/register", cfg.AuthHandler.Register)
	auth.Post("/login", cfg.AuthHandler.Login)

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
	orders.Post("/", cfg.OrderHandler.CreateOrder)
	orders.Post("/checkout", cfg.OrderHandler.CreateOrder)
	orders.Get("/", cfg.OrderHandler.GetUserOrders)
	orders.Get("/:code", cfg.OrderHandler.GetOrderByCode)

	// Payments
	api.Post("/payments/:order_code/snap", userAuth, cfg.PaymentHandler.CreateSnapToken)

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

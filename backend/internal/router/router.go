package router

import (
	"time"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/config"
	"github.com/MrVirul/icecream-inventory/internal/handler"
	"github.com/MrVirul/icecream-inventory/internal/middleware"
	"github.com/MrVirul/icecream-inventory/internal/repository"
	"github.com/MrVirul/icecream-inventory/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit:    1 << 20,
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	// Repos
	userRepo := repository.NewUserRepository(db)
	revokedRepo := repository.NewRevokedTokenRepository(db)
	itemRepo := repository.NewItemRepository(db)
	customerRepo := repository.NewCustomerRepository(db)

	// Services
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTTTLMinutes)
	authSvc := service.NewAuthService(userRepo, revokedRepo, jwtSvc, cfg.BcryptCost)
	itemSvc := service.NewItemService(itemRepo)
	customerSvc := service.NewCustomerService(customerRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	healthH := handler.NewHealthHandler(db)
	itemH := handler.NewItemHandler(itemSvc)
	customerH := handler.NewCustomerHandler(customerSvc)

	// Middleware
	authMw := middleware.NewAuth(jwtSvc, db)

	// Routes
	api := app.Group("/api")

	// Login rate limit
	loginLimiter := limiter.New(limiter.Config{
		Max:        cfg.LoginRateLimitPerMin,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many requests")
		},
	})

	api.Post("/auth/login", loginLimiter, authH.Login)

	protected := api.Group("", authMw.Authenticate())
	protected.Post("/auth/logout", authH.Logout)

	// Items
	items := protected.Group("/items")
	items.Get("/", middleware.RequirePermission(auth.ItemsRead), itemH.List)
	items.Get("/:id", middleware.RequirePermission(auth.ItemsRead), itemH.Get)
	items.Post("/", middleware.RequirePermission(auth.ItemsCreate), itemH.Create)
	items.Put("/:id", middleware.RequirePermission(auth.ItemsUpdate), itemH.Replace)
	items.Patch("/:id", middleware.RequirePermission(auth.ItemsUpdate), itemH.Patch)
	items.Delete("/:id", middleware.RequirePermission(auth.ItemsDelete), itemH.Delete)

	// Customers
	customers := protected.Group("/customers")
	customers.Get("/", middleware.RequirePermission(auth.CustomersRead), customerH.List)
	customers.Get("/:id", middleware.RequirePermission(auth.CustomersRead), customerH.Get)
	customers.Post("/", middleware.RequirePermission(auth.CustomersCreate), customerH.Create)
	customers.Put("/:id", middleware.RequirePermission(auth.CustomersUpdate), customerH.Replace)
	customers.Patch("/:id", middleware.RequirePermission(auth.CustomersUpdate), customerH.Patch)
	customers.Delete("/:id", middleware.RequirePermission(auth.CustomersDelete), customerH.Delete)

	app.Get("/healthz", healthH.Healthz)


	return app
}

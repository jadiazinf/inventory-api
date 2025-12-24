package api

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/jadiazinf/inventory/internal/adapters/http/handlers"
	"github.com/jadiazinf/inventory/internal/adapters/http/middleware"
	"github.com/jadiazinf/inventory/internal/config"
	"github.com/jadiazinf/inventory/internal/platform/i18n"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	ProductHandler     *handlers.ProductHandler
	CustomerHandler    *handlers.CustomerHandler
	SaleHandler        *handlers.SaleHandler
	ReservationHandler *handlers.ReservationHandler
	InventoryHandler   *handlers.InventoryHandler
}

type Server struct {
	app            *fiber.App
	cfg            *config.Config
	db             *gorm.DB
	handlers       *Handlers
	authMiddleware *middleware.AuthMiddleware
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName:     "Inventory API",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New())
	app.Use(i18n.Middleware())

	s := &Server{
		app: app,
		cfg: cfg,
		db:  db,
	}

	return s
}

// SetHandlers sets the HTTP handlers for the server and sets up routes
func (s *Server) SetHandlers(handlers *Handlers) {
	s.handlers = handlers
	s.setupRoutes() // Setup routes after handlers are set
}

// SetAuthMiddleware sets the authentication middleware
func (s *Server) SetAuthMiddleware(authMiddleware *middleware.AuthMiddleware) {
	s.authMiddleware = authMiddleware
}

func (s *Server) Run() error {
	log.Info("Server listening on port ", s.cfg.AppPort)
	return s.app.Listen(":" + s.cfg.AppPort)
}

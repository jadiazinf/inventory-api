package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jadiazinf/inventory/internal/adapters/http/handlers"
	"github.com/jadiazinf/inventory/internal/adapters/http/middleware"
	postgresRepo "github.com/jadiazinf/inventory/internal/adapters/repository/postgres"
	"github.com/jadiazinf/inventory/internal/api"
	"github.com/jadiazinf/inventory/internal/config"
	"github.com/jadiazinf/inventory/internal/core/services"
	"github.com/jadiazinf/inventory/internal/platform/database"
	"github.com/jadiazinf/inventory/internal/platform/firebase"
	"github.com/jadiazinf/inventory/internal/platform/i18n"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger
	log.Info("Starting Inventory API...")

	// 3. Initialize i18n
	i18n.InitI18n()

	// 4. Initialize Database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		os.Exit(1)
	} else {
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		log.Info("Database connected successfully")
	}

	// 5. Initialize Firebase
	firebaseApp, err := firebase.NewFirebaseApp(cfg)
	if err != nil {
		log.Error("Failed to initialize Firebase: ", err)
		// Don't exit - some features may work without Firebase
	} else {
		log.Info("Firebase initialized successfully")
	}

	// 6. Initialize Repositories
	log.Info("Initializing repositories...")
	userRepo := postgresRepo.NewUserRepository(db)
	productRepo := postgresRepo.NewProductRepository(db)
	customerRepo := postgresRepo.NewCustomerRepository(db)
	customerChildRepo := postgresRepo.NewCustomerChildRepository(db)
	inventoryRepo := postgresRepo.NewInventoryRepository(db)
	saleRepo := postgresRepo.NewSaleRepository(db)
	reservationRepo := postgresRepo.NewReservationRepository(db)
	arRepo := postgresRepo.NewAccountsReceivableRepository(db)

	// 7. Initialize Services
	log.Info("Initializing services...")
	productService := services.NewProductService(productRepo, inventoryRepo, db)
	inventoryService := services.NewInventoryService(inventoryRepo, productRepo, db)
	arService := services.NewAccountsReceivableService(arRepo, db)
	notificationService := services.NewNotificationService(reservationRepo, customerRepo, db)
	saleService := services.NewSaleService(saleRepo, productRepo, inventoryRepo, customerRepo, db)
	reservationService := services.NewReservationService(
		reservationRepo,
		customerRepo,
		productRepo,
		inventoryRepo,
		saleRepo,
		notificationService,
		db,
	)

	// 8. Initialize Middleware
	log.Info("Initializing middleware...")
	var authMiddleware *middleware.AuthMiddleware
	if firebaseApp != nil {
		authMiddleware = middleware.NewAuthMiddleware(userRepo)
	} else {
		log.Warn("Running without authentication middleware (Firebase not initialized)")
	}

	// 9. Initialize Handlers
	log.Info("Initializing handlers...")
	apiHandlers := &api.Handlers{
		ProductHandler:     handlers.NewProductHandler(productService),
		CustomerHandler:    handlers.NewCustomerHandler(customerRepo, customerChildRepo),
		SaleHandler:        handlers.NewSaleHandler(saleService, arService),
		ReservationHandler: handlers.NewReservationHandler(reservationService),
		InventoryHandler:   handlers.NewInventoryHandler(inventoryService),
	}

	log.Info("All handlers initialized successfully")

	// 10. Create Server
	log.Info("Creating server...")
	server := api.NewServer(cfg, db)
	server.SetHandlers(apiHandlers)
	if authMiddleware != nil {
		server.SetAuthMiddleware(authMiddleware)
	}

	// 11. Start Server
	log.Info("Starting server on port ", cfg.AppPort)
	if err := server.Run(); err != nil {
		log.Error("Failed to start server: ", err)
		os.Exit(1)
	}
}

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jadiazinf/inventory/internal/health"
)

func (s *Server) setupRoutes() {
	// Health check (public)
	s.app.Get("/health", health.NewHandler().Check)

	// API v1 group
	api := s.app.Group("/api/v1")

	// Health routes
	healthHandler := health.NewHandler()
	api.Get("/health", healthHandler.Check)
	api.Get("/greet/:name", healthHandler.Greet)

	// Setup resource routes if handlers are available
	if s.handlers != nil {
		s.setupProductRoutes(api)
		s.setupCustomerRoutes(api)
		s.setupSaleRoutes(api)
		s.setupReservationRoutes(api)
		s.setupInventoryRoutes(api)
	}
}

func (s *Server) setupProductRoutes(api fiber.Router) {
	if s.handlers.ProductHandler == nil {
		return
	}

	products := api.Group("/products")

	// Public routes
	products.Get("/", s.handlers.ProductHandler.ListProducts)
	products.Get("/search", s.handlers.ProductHandler.SearchProducts)
	products.Get("/:id", s.handlers.ProductHandler.GetProduct)
	products.Get("/sku/:sku", s.handlers.ProductHandler.GetProductBySKU)

	// Protected routes (require authentication)
	if s.authMiddleware != nil {
		products.Post("/", s.authMiddleware.Authenticate(), s.handlers.ProductHandler.CreateProduct)
		products.Put("/:id", s.authMiddleware.Authenticate(), s.handlers.ProductHandler.UpdateProduct)
		products.Put("/:id/price", s.authMiddleware.Authenticate(), s.handlers.ProductHandler.UpdatePrice)
		products.Delete("/:id", s.authMiddleware.Authenticate(), s.handlers.ProductHandler.DeleteProduct)
	}
}

func (s *Server) setupCustomerRoutes(api fiber.Router) {
	if s.handlers.CustomerHandler == nil {
		return
	}

	customers := api.Group("/customers")

	// All customer routes require authentication
	if s.authMiddleware != nil {
		customers.Use(s.authMiddleware.Authenticate())
	}

	customers.Get("/", s.handlers.CustomerHandler.ListCustomers)
	customers.Get("/:id", s.handlers.CustomerHandler.GetCustomer)
	customers.Get("/:id/with-children", s.handlers.CustomerHandler.GetCustomerWithChildren)
	customers.Get("/tax-id/:taxId", s.handlers.CustomerHandler.GetCustomerByTaxID)
	customers.Post("/", s.handlers.CustomerHandler.CreateCustomer)
	customers.Put("/:id", s.handlers.CustomerHandler.UpdateCustomer)
	customers.Delete("/:id", s.handlers.CustomerHandler.DeleteCustomer)

	// Children management
	customers.Post("/:id/children", s.handlers.CustomerHandler.AddChild)
	customers.Get("/:id/children", s.handlers.CustomerHandler.GetChildren)
	customers.Put("/:id/children/:childId", s.handlers.CustomerHandler.UpdateChild)
	customers.Delete("/:id/children/:childId", s.handlers.CustomerHandler.DeleteChild)
	customers.Put("/:id/loyalty-points", s.handlers.CustomerHandler.UpdateLoyaltyPoints)
}

func (s *Server) setupSaleRoutes(api fiber.Router) {
	if s.handlers.SaleHandler == nil {
		return
	}

	sales := api.Group("/sales")

	// All sales routes require authentication
	if s.authMiddleware != nil {
		sales.Use(s.authMiddleware.Authenticate())
	}

	sales.Get("/", s.handlers.SaleHandler.ListSales)
	sales.Get("/daily", s.handlers.SaleHandler.GetDailySales)
	sales.Get("/:id", s.handlers.SaleHandler.GetSale)
	sales.Get("/invoice/:invoice", s.handlers.SaleHandler.GetSaleByInvoiceNumber)
	sales.Post("/", s.handlers.SaleHandler.CreateSale)
	sales.Post("/credit", s.handlers.SaleHandler.CreateCreditSale)
	sales.Post("/:id/cancel", s.handlers.SaleHandler.CancelSale)

	// Accounts Receivable routes
	ar := api.Group("/accounts-receivable")
	if s.authMiddleware != nil {
		ar.Use(s.authMiddleware.Authenticate())
	}
	ar.Get("/:id", s.handlers.SaleHandler.GetAccountsReceivable)
	ar.Post("/:id/payments", s.handlers.SaleHandler.RegisterPayment)
}

func (s *Server) setupReservationRoutes(api fiber.Router) {
	if s.handlers.ReservationHandler == nil {
		return
	}

	reservations := api.Group("/reservations")

	// All reservation routes require authentication
	if s.authMiddleware != nil {
		reservations.Use(s.authMiddleware.Authenticate())
	}

	reservations.Get("/", s.handlers.ReservationHandler.ListReservations)
	reservations.Get("/:id", s.handlers.ReservationHandler.GetReservation)
	reservations.Get("/number/:number", s.handlers.ReservationHandler.GetReservationByNumber)
	reservations.Post("/", s.handlers.ReservationHandler.CreateReservation)
	reservations.Post("/:id/confirm", s.handlers.ReservationHandler.ConfirmReservation)
	reservations.Post("/:id/fulfill", s.handlers.ReservationHandler.FulfillReservation)
	reservations.Post("/:id/cancel", s.handlers.ReservationHandler.CancelReservation)
}

func (s *Server) setupInventoryRoutes(api fiber.Router) {
	if s.handlers.InventoryHandler == nil {
		return
	}

	inventory := api.Group("/inventory")

	// All inventory routes require authentication
	if s.authMiddleware != nil {
		inventory.Use(s.authMiddleware.Authenticate())
	}

	inventory.Get("/product/:productId/warehouse/:warehouseId", s.handlers.InventoryHandler.GetInventory)
	inventory.Get("/warehouse/:warehouseId", s.handlers.InventoryHandler.GetWarehouseInventory)
	inventory.Get("/product/:productId", s.handlers.InventoryHandler.GetProductInventory)
	inventory.Get("/check-availability", s.handlers.InventoryHandler.CheckAvailability)

	// Movement operations
	inventory.Post("/movements/inbound", s.handlers.InventoryHandler.RegisterInboundMovement)
	inventory.Post("/movements/outbound", s.handlers.InventoryHandler.RegisterOutboundMovement)
	inventory.Post("/movements/adjustment", s.handlers.InventoryHandler.RegisterAdjustment)
	inventory.Get("/movements/product/:productId", s.handlers.InventoryHandler.GetProductMovements)
	inventory.Get("/movements/warehouse/:warehouseId", s.handlers.InventoryHandler.GetWarehouseMovements)
}

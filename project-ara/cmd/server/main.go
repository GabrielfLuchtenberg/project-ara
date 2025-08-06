package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"project-ara/internal/database"
	"project-ara/internal/handlers"
	"project-ara/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	// Set up logger
	logrus.SetLevel(logrus.InfoLevel)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	whatsappService := services.NewWhatsAppService()
	nlpService := services.NewNLPService()
	voiceService := services.NewVoiceService()
	ocrService := services.NewOCRService()
	transactionService := services.NewTransactionService(db)
	userService := services.NewUserService(db)

	// Initialize Phase 3 services
	reportingService := services.NewFinancialReportingService(transactionService, userService)
	subscriptionService := services.NewSubscriptionService(userService, transactionService, reportingService)

	// Initialize handlers
	whatsappHandler := handlers.NewWhatsAppHandler(whatsappService, nlpService, voiceService, ocrService, transactionService, userService, reportingService, subscriptionService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize Phase 3 handlers
	financialHandler := handlers.NewFinancialHandler(transactionService, reportingService, subscriptionService)

	// Set up router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// WhatsApp webhook
		api.POST("/webhook/whatsapp", whatsappHandler.HandleWebhook)

		// Transaction endpoints
		api.POST("/transactions", whatsappHandler.CreateTransaction)
		api.GET("/users/:id/summary", whatsappHandler.GetUserSummary)

		// Phase 3: Financial endpoints
		financial := api.Group("/financial")
		{
			financial.GET("/users/:userID/summary", financialHandler.GetFinancialSummary)
			financial.GET("/users/:userID/report", financialHandler.GetDetailedReport)
			financial.GET("/users/:userID/balance", financialHandler.GetUserBalance)
			financial.GET("/users/:userID/transactions", financialHandler.GetUserTransactions)
			financial.PUT("/transactions/:transactionID/correct", financialHandler.CorrectTransaction)
			financial.GET("/users/:userID/categories", financialHandler.GetTopCategories)
		}

		// Phase 3: Subscription endpoints
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.GET("/users/:userID/trial-status", financialHandler.GetTrialStatus)
			subscriptions.GET("/users/:userID/info", financialHandler.GetSubscriptionInfo)
			subscriptions.POST("/users/:userID", financialHandler.CreateSubscription)
			subscriptions.DELETE("/users/:userID", financialHandler.CancelSubscription)
			subscriptions.POST("/webhook/payment", financialHandler.ProcessPaymentWebhook)
		}

		// Legacy endpoints (for backward compatibility)
		api.POST("/subscriptions", whatsappHandler.CreateSubscription)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

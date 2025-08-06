package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-ara/internal/services"
)

type FinancialHandler struct {
	transactionService  *services.TransactionService
	reportingService    *services.FinancialReportingService
	subscriptionService *services.SubscriptionService
}

func NewFinancialHandler(transactionService *services.TransactionService, reportingService *services.FinancialReportingService, subscriptionService *services.SubscriptionService) *FinancialHandler {
	return &FinancialHandler{
		transactionService:  transactionService,
		reportingService:    reportingService,
		subscriptionService: subscriptionService,
	}
}

// GetFinancialSummary returns a conversational financial summary
func (h *FinancialHandler) GetFinancialSummary(c *gin.Context) {
	userID := c.Param("userID")
	period := c.DefaultQuery("period", "today")

	summary, err := h.reportingService.GenerateConversationalSummary(userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate financial summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"period":  period,
		"user_id": userID,
	})
}

// GetDetailedReport returns a comprehensive financial report
func (h *FinancialHandler) GetDetailedReport(c *gin.Context) {
	userID := c.Param("userID")
	period := c.DefaultQuery("period", "today")

	report, err := h.reportingService.GenerateDetailedReport(userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate detailed report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetUserBalance returns the current user balance
func (h *FinancialHandler) GetUserBalance(c *gin.Context) {
	userID := c.Param("userID")

	balance, err := h.transactionService.GetUserBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user balance",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"balance":  balance,
		"currency": "BRL",
	})
}

// GetUserTransactions returns user's transaction history
func (h *FinancialHandler) GetUserTransactions(c *gin.Context) {
	userID := c.Param("userID")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// CorrectTransaction allows users to correct a transaction
func (h *FinancialHandler) CorrectTransaction(c *gin.Context) {
	transactionID := c.Param("transactionID")

	var request struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	transaction, err := h.transactionService.CorrectTransaction(transactionID, request.Amount, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to correct transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction corrected successfully",
		"transaction": transaction,
	})
}

// GetTopCategories returns top spending/income categories
func (h *FinancialHandler) GetTopCategories(c *gin.Context) {
	userID := c.Param("userID")
	period := c.DefaultQuery("period", "today")
	limitStr := c.DefaultQuery("limit", "5")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	categories, err := h.transactionService.GetTopCategories(userID, period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get top categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"period":     period,
		"categories": categories,
	})
}

// GetTrialStatus returns user's trial status
func (h *FinancialHandler) GetTrialStatus(c *gin.Context) {
	userID := c.Param("userID")

	status, err := h.subscriptionService.CheckTrialStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get trial status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetSubscriptionInfo returns detailed subscription information
func (h *FinancialHandler) GetSubscriptionInfo(c *gin.Context) {
	userID := c.Param("userID")

	info, err := h.subscriptionService.GetSubscriptionInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get subscription info",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// CreateSubscription creates a new subscription
func (h *FinancialHandler) CreateSubscription(c *gin.Context) {
	userID := c.Param("userID")

	var request struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	subscription, err := h.subscriptionService.CreateSubscription(userID, request.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create subscription",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Subscription created successfully",
		"subscription": subscription,
	})
}

// CancelSubscription cancels the user's subscription
func (h *FinancialHandler) CancelSubscription(c *gin.Context) {
	userID := c.Param("userID")

	if err := h.subscriptionService.CancelSubscription(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel subscription",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription cancelled successfully",
	})
}

// ProcessPaymentWebhook handles payment gateway webhooks
func (h *FinancialHandler) ProcessPaymentWebhook(c *gin.Context) {
	var webhookData map[string]interface{}

	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid webhook data",
			"details": err.Error(),
		})
		return
	}

	if err := h.subscriptionService.ProcessPaymentWebhook(webhookData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process payment webhook",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment webhook processed successfully",
	})
}

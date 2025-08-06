package services

import (
	"fmt"
	"time"

	"project-ara/internal/models"
)

type SubscriptionService struct {
	userService        *UserService
	transactionService *TransactionService
	reportingService   *FinancialReportingService
}

func NewSubscriptionService(userService *UserService, transactionService *TransactionService, reportingService *FinancialReportingService) *SubscriptionService {
	return &SubscriptionService{
		userService:        userService,
		transactionService: transactionService,
		reportingService:   reportingService,
	}
}

// CheckTrialStatus checks if user should be prompted for subscription
func (s *SubscriptionService) CheckTrialStatus(userID string) (*TrialStatus, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	status := &TrialStatus{
		UserID:                      userID,
		SubscriptionStatus:          user.SubscriptionStatus,
		TrialTransactionsCount:      user.TrialTransactionsCount,
		RemainingTrialTransactions:  50 - user.TrialTransactionsCount,
		IsTrialExpired:              user.IsTrialExpired(),
		ShouldPromptForSubscription: false,
	}

	// Determine if we should prompt for subscription
	if user.SubscriptionStatus == "trial" {
		if user.TrialTransactionsCount >= 45 { // Prompt when 5 transactions remaining
			status.ShouldPromptForSubscription = true
		}
		if user.TrialTransactionsCount >= 50 { // Trial expired
			status.IsTrialExpired = true
			status.ShouldPromptForSubscription = true
		}
	}

	return status, nil
}

// CreateSubscription creates a new subscription for the user
func (s *SubscriptionService) CreateSubscription(userID string, paymentMethod string) (*Subscription, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user already has active subscription
	if user.SubscriptionStatus == "active" {
		return nil, fmt.Errorf("user already has active subscription")
	}

	// Calculate subscription expiry (30 days from now)
	expiresAt := time.Now().AddDate(0, 0, 30)

	// Update user subscription status
	if err := s.userService.UpdateSubscriptionStatus(userID, "active"); err != nil {
		return nil, fmt.Errorf("failed to update subscription status: %w", err)
	}

	// Update subscription expiry
	if err := s.updateSubscriptionExpiry(userID, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to update subscription expiry: %w", err)
	}

	subscription := &Subscription{
		UserID:        userID,
		Status:        "active",
		PaymentMethod: paymentMethod,
		Amount:        9.90, // R$ 9,90/month
		Currency:      "BRL",
		CreatedAt:     time.Now(),
		ExpiresAt:     expiresAt,
	}

	return subscription, nil
}

// CancelSubscription cancels the user's subscription
func (s *SubscriptionService) CancelSubscription(userID string) error {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.SubscriptionStatus != "active" {
		return fmt.Errorf("user does not have active subscription")
	}

	// Update subscription status to cancelled
	if err := s.userService.UpdateSubscriptionStatus(userID, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// RenewSubscription renews the user's subscription
func (s *SubscriptionService) RenewSubscription(userID string) error {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.SubscriptionStatus != "active" {
		return fmt.Errorf("user does not have active subscription")
	}

	// Extend subscription by 30 days
	newExpiry := time.Now().AddDate(0, 0, 30)
	if user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.After(time.Now()) {
		// If subscription is still active, extend from current expiry
		newExpiry = user.SubscriptionExpiresAt.AddDate(0, 0, 30)
	}

	if err := s.updateSubscriptionExpiry(userID, newExpiry); err != nil {
		return fmt.Errorf("failed to renew subscription: %w", err)
	}

	return nil
}

// GetSubscriptionInfo returns detailed subscription information
func (s *SubscriptionService) GetSubscriptionInfo(userID string) (*SubscriptionInfo, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get trial status
	trialStatus, err := s.CheckTrialStatus(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trial status: %w", err)
	}

	info := &SubscriptionInfo{
		UserID:                     userID,
		SubscriptionStatus:         user.SubscriptionStatus,
		TrialTransactionsCount:     user.TrialTransactionsCount,
		RemainingTrialTransactions: trialStatus.RemainingTrialTransactions,
		IsTrialExpired:             trialStatus.IsTrialExpired,
		SubscriptionExpiresAt:      user.SubscriptionExpiresAt,
		MonthlyPrice:               9.90,
		Currency:                   "BRL",
	}

	// Calculate days until expiry
	if user.SubscriptionExpiresAt != nil {
		info.DaysUntilExpiry = int(user.SubscriptionExpiresAt.Sub(time.Now()).Hours() / 24)
	}

	return info, nil
}

// ProcessPaymentWebhook handles payment gateway webhooks
func (s *SubscriptionService) ProcessPaymentWebhook(webhookData map[string]interface{}) error {
	// Extract webhook data
	userID, ok := webhookData["user_id"].(string)
	if !ok {
		return fmt.Errorf("invalid user_id in webhook data")
	}

	paymentStatus, ok := webhookData["status"].(string)
	if !ok {
		return fmt.Errorf("invalid payment status in webhook data")
	}

	_, ok = webhookData["payment_id"].(string)
	if !ok {
		return fmt.Errorf("invalid payment_id in webhook data")
	}

	switch paymentStatus {
	case "approved":
		// Payment successful - activate subscription
		if err := s.userService.UpdateSubscriptionStatus(userID, "active"); err != nil {
			return fmt.Errorf("failed to activate subscription: %w", err)
		}

		// Set subscription expiry
		expiresAt := time.Now().AddDate(0, 0, 30)
		if err := s.updateSubscriptionExpiry(userID, expiresAt); err != nil {
			return fmt.Errorf("failed to set subscription expiry: %w", err)
		}

	case "failed", "cancelled":
		// Payment failed - keep user in trial or cancelled status
		if err := s.userService.UpdateSubscriptionStatus(userID, "trial"); err != nil {
			return fmt.Errorf("failed to revert subscription status: %w", err)
		}

	case "refunded":
		// Payment refunded - cancel subscription
		if err := s.userService.UpdateSubscriptionStatus(userID, "cancelled"); err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}
	}

	return nil
}

// updateSubscriptionExpiry updates the subscription expiry date
func (s *SubscriptionService) updateSubscriptionExpiry(userID string, expiresAt time.Time) error {
	return s.userService.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("subscription_expires_at", expiresAt).
		Error
}

type TrialStatus struct {
	UserID                      string `json:"user_id"`
	SubscriptionStatus          string `json:"subscription_status"`
	TrialTransactionsCount      int    `json:"trial_transactions_count"`
	RemainingTrialTransactions  int    `json:"remaining_trial_transactions"`
	IsTrialExpired              bool   `json:"is_trial_expired"`
	ShouldPromptForSubscription bool   `json:"should_prompt_for_subscription"`
}

type Subscription struct {
	UserID        string    `json:"user_id"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type SubscriptionInfo struct {
	UserID                     string     `json:"user_id"`
	SubscriptionStatus         string     `json:"subscription_status"`
	TrialTransactionsCount     int        `json:"trial_transactions_count"`
	RemainingTrialTransactions int        `json:"remaining_trial_transactions"`
	IsTrialExpired             bool       `json:"is_trial_expired"`
	SubscriptionExpiresAt      *time.Time `json:"subscription_expires_at,omitempty"`
	DaysUntilExpiry            int        `json:"days_until_expiry"`
	MonthlyPrice               float64    `json:"monthly_price"`
	Currency                   string     `json:"currency"`
}

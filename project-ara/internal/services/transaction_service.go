package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"project-ara/internal/models"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) CreateTransaction(userID string, amount float64, description string, transactionType models.TransactionType, source models.TransactionSource) (*models.Transaction, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	transaction := &models.Transaction{
		UserID:          userUUID,
		Amount:          amount,
		Description:     description,
		TransactionType: transactionType,
		Source:          source,
		CreatedAt:       time.Now(),
	}

	if err := s.db.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user's trial transaction count
	if err := s.updateUserTrialCount(userUUID.String()); err != nil {
		return nil, fmt.Errorf("failed to update trial count: %w", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetUserTransactions(userID string, limit int) ([]models.Transaction, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var transactions []models.Transaction
	if err := s.db.Where("user_id = ?", userUUID).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (s *TransactionService) GetFinancialSummary(userID string) (*FinancialSummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get today's transactions
	today := time.Now().Truncate(24 * time.Hour)
	var todayTransactions []models.Transaction
	if err := s.db.Where("user_id = ? AND created_at >= ?", userUUID, today).Find(&todayTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get today's transactions: %w", err)
	}

	// Calculate summary
	summary := &FinancialSummary{
		UserID: userID,
		Date:   today,
	}

	for _, t := range todayTransactions {
		if t.IsIncome() {
			summary.TotalIncome += t.Amount
		} else {
			summary.TotalExpenses += t.Amount
		}
	}

	summary.Profit = summary.TotalIncome - summary.TotalExpenses

	return summary, nil
}

// New Phase 3 methods

func (s *TransactionService) GetPeriodSummary(userID string, period string) (*PeriodSummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = now.Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var transactions []models.Transaction
	if err := s.db.Where("user_id = ? AND created_at >= ? AND created_at < ?",
		userUUID, startDate, endDate).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get period transactions: %w", err)
	}

	summary := &PeriodSummary{
		UserID:    userID,
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	}

	for _, t := range transactions {
		if t.IsIncome() {
			summary.TotalIncome += t.Amount
		} else {
			summary.TotalExpenses += t.Amount
		}
	}

	summary.Profit = summary.TotalIncome - summary.TotalExpenses
	summary.TransactionCount = len(transactions)

	return summary, nil
}

func (s *TransactionService) CorrectTransaction(transactionID string, correctedAmount float64, correctedDescription string) (*models.Transaction, error) {
	transactionUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	var transaction models.Transaction
	if err := s.db.Where("id = ?", transactionUUID).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Store correction data
	correctionData := map[string]interface{}{
		"original_amount":       transaction.Amount,
		"original_description":  transaction.Description,
		"corrected_amount":      correctedAmount,
		"corrected_description": correctedDescription,
		"corrected_at":          time.Now(),
	}

	// Update transaction
	updates := map[string]interface{}{
		"amount":          correctedAmount,
		"description":     correctedDescription,
		"corrected_at":    time.Now(),
		"correction_data": correctionData,
	}

	if err := s.db.Model(&transaction).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to correct transaction: %w", err)
	}

	return &transaction, nil
}

func (s *TransactionService) GetTransactionByID(transactionID string) (*models.Transaction, error) {
	transactionUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	var transaction models.Transaction
	if err := s.db.Where("id = ?", transactionUUID).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return &transaction, nil
}

func (s *TransactionService) GetUserBalance(userID string) (float64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}

	var result struct {
		Balance float64
	}

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END), 0) as balance
		FROM transactions 
		WHERE user_id = ?
	`

	if err := s.db.Raw(query, userUUID).Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return result.Balance, nil
}

func (s *TransactionService) GetTopCategories(userID string, period string, limit int) ([]CategorySummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var startDate time.Time
	now := time.Now()

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.Truncate(24 * time.Hour) // today
	}

	var results []CategorySummary
	query := `
		SELECT 
			description,
			transaction_type,
			COUNT(*) as count,
			SUM(amount) as total_amount
		FROM transactions 
		WHERE user_id = ? AND created_at >= ?
		GROUP BY description, transaction_type
		ORDER BY total_amount DESC
		LIMIT ?
	`

	if err := s.db.Raw(query, userUUID, startDate, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}

	return results, nil
}

func (s *TransactionService) updateUserTrialCount(userID string) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("trial_transactions_count", gorm.Expr("trial_transactions_count + 1")).
		Error
}

type FinancialSummary struct {
	UserID        string    `json:"user_id"`
	Date          time.Time `json:"date"`
	TotalIncome   float64   `json:"total_income"`
	TotalExpenses float64   `json:"total_expenses"`
	Profit        float64   `json:"profit"`
}

type PeriodSummary struct {
	UserID           string    `json:"user_id"`
	Period           string    `json:"period"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	TotalIncome      float64   `json:"total_income"`
	TotalExpenses    float64   `json:"total_expenses"`
	Profit           float64   `json:"profit"`
	TransactionCount int       `json:"transaction_count"`
}

type CategorySummary struct {
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type"`
	Count           int     `json:"count"`
	TotalAmount     float64 `json:"total_amount"`
}

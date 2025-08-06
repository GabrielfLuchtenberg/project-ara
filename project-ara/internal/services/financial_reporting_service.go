package services

import (
	"fmt"
	"strings"
	"time"

	"project-ara/internal/models"
)

type FinancialReportingService struct {
	transactionService *TransactionService
	userService        *UserService
}

func NewFinancialReportingService(transactionService *TransactionService, userService *UserService) *FinancialReportingService {
	return &FinancialReportingService{
		transactionService: transactionService,
		userService:        userService,
	}
}

// GenerateConversationalSummary creates a user-friendly financial summary in Portuguese
func (s *FinancialReportingService) GenerateConversationalSummary(userID string, period string) (string, error) {
	summary, err := s.transactionService.GetPeriodSummary(userID, period)
	if err != nil {
		return "", fmt.Errorf("failed to get period summary: %w", err)
	}

	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return s.formatConversationalSummary(summary, user, period), nil
}

// GenerateDetailedReport creates a comprehensive financial report
func (s *FinancialReportingService) GenerateDetailedReport(userID string, period string) (*DetailedReport, error) {
	summary, err := s.transactionService.GetPeriodSummary(userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get period summary: %w", err)
	}

	balance, err := s.transactionService.GetUserBalance(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	topCategories, err := s.transactionService.GetTopCategories(userID, period, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}

	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	report := &DetailedReport{
		UserID:         userID,
		Period:         period,
		Summary:        summary,
		CurrentBalance: balance,
		TopCategories:  topCategories,
		User:           user,
		GeneratedAt:    time.Now(),
	}

	// Calculate trends
	report.CalculateTrends()

	return report, nil
}

// GenerateTrialStatusMessage creates a message about user's trial status
func (s *FinancialReportingService) GenerateTrialStatusMessage(userID string) (string, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	remainingTransactions := 50 - user.TrialTransactionsCount

	if user.SubscriptionStatus == "active" {
		return "✅ Sua assinatura está ativa! Você pode registrar transações ilimitadas.", nil
	}

	if remainingTransactions <= 0 {
		return "⚠️ Você atingiu o limite de 50 transações do período de teste. Para continuar usando o Ara, assine o plano premium por apenas R$ 9,90/mês e tenha transações ilimitadas! 💰", nil
	}

	if remainingTransactions <= 10 {
		return fmt.Sprintf("⚠️ Você tem apenas %d transações restantes no período de teste. Considere assinar o plano premium por R$ 9,90/mês para transações ilimitadas! 💰", remainingTransactions), nil
	}

	return fmt.Sprintf("📊 Você tem %d transações restantes no período de teste.", remainingTransactions), nil
}

// GenerateConversionMessage creates a compelling conversion message
func (s *FinancialReportingService) GenerateConversionMessage(userID string) (string, error) {
	_, err := s.userService.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Get user's financial performance
	todaySummary, err := s.transactionService.GetPeriodSummary(userID, "today")
	if err != nil {
		return "", fmt.Errorf("failed to get today's summary: %w", err)
	}

	weekSummary, err := s.transactionService.GetPeriodSummary(userID, "week")
	if err != nil {
		return "", fmt.Errorf("failed to get week's summary: %w", err)
	}

	balance, err := s.transactionService.GetUserBalance(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	message := "🚀 **Ara Premium - Transforme seu negócio!**\n\n"
	message += "Veja como o Ara está ajudando você:\n\n"

	if todaySummary.Profit > 0 {
		message += fmt.Sprintf("💰 Hoje: R$ %.2f de lucro\n", todaySummary.Profit)
	}

	if weekSummary.Profit > 0 {
		message += fmt.Sprintf("📈 Esta semana: R$ %.2f de lucro\n", weekSummary.Profit)
	}

	message += fmt.Sprintf("💳 Saldo atual: R$ %.2f\n\n", balance)

	message += "**Benefícios Premium:**\n"
	message += "✅ Transações ilimitadas\n"
	message += "✅ Relatórios avançados\n"
	message += "✅ Categorização automática\n"
	message += "✅ Backup na nuvem\n"
	message += "✅ Suporte prioritário\n\n"

	message += "💎 **Apenas R$ 9,90/mês**\n"
	message += "Menos que um café por dia! ☕\n\n"

	message += "Para assinar, responda: *ASSINAR*"

	return message, nil
}

// formatConversationalSummary formats the summary in conversational Portuguese
func (s *FinancialReportingService) formatConversationalSummary(summary *PeriodSummary, user *models.User, period string) string {
	var message strings.Builder

	// Period-specific greeting
	switch period {
	case "today":
		message.WriteString("📊 **Resumo de hoje:**\n\n")
	case "week":
		message.WriteString("📈 **Resumo da semana:**\n\n")
	case "month":
		message.WriteString("📅 **Resumo do mês:**\n\n")
	default:
		message.WriteString("📊 **Resumo financeiro:**\n\n")
	}

	// Income section
	if summary.TotalIncome > 0 {
		message.WriteString(fmt.Sprintf("💰 **Receitas:** R$ %.2f\n", summary.TotalIncome))
	}

	// Expenses section
	if summary.TotalExpenses > 0 {
		message.WriteString(fmt.Sprintf("💸 **Despesas:** R$ %.2f\n", summary.TotalExpenses))
	}

	// Profit/Loss section
	message.WriteString("\n")
	if summary.Profit > 0 {
		message.WriteString(fmt.Sprintf("✅ **Lucro:** R$ %.2f\n", summary.Profit))
	} else if summary.Profit < 0 {
		message.WriteString(fmt.Sprintf("❌ **Prejuízo:** R$ %.2f\n", -summary.Profit))
	} else {
		message.WriteString("⚖️ **Empate:** R$ 0,00\n")
	}

	// Transaction count
	message.WriteString(fmt.Sprintf("\n📝 **Total de transações:** %d\n", summary.TransactionCount))

	// Trial status
	if user.SubscriptionStatus == "trial" {
		remaining := 50 - user.TrialTransactionsCount
		message.WriteString(fmt.Sprintf("\n🎯 **Transações restantes no teste:** %d\n", remaining))
	}

	return message.String()
}

type DetailedReport struct {
	UserID         string            `json:"user_id"`
	Period         string            `json:"period"`
	Summary        *PeriodSummary    `json:"summary"`
	CurrentBalance float64           `json:"current_balance"`
	TopCategories  []CategorySummary `json:"top_categories"`
	User           *models.User      `json:"user"`
	GeneratedAt    time.Time         `json:"generated_at"`
	Trends         *TrendAnalysis    `json:"trends,omitempty"`
}

type TrendAnalysis struct {
	ProfitTrend     string   `json:"profit_trend"` // "increasing", "decreasing", "stable"
	IncomeTrend     string   `json:"income_trend"`
	ExpenseTrend    string   `json:"expense_trend"`
	GrowthRate      float64  `json:"growth_rate"`
	Recommendations []string `json:"recommendations"`
}

func (r *DetailedReport) CalculateTrends() {
	// This would compare with previous periods to calculate trends
	// For now, we'll set basic trends based on current data
	r.Trends = &TrendAnalysis{
		ProfitTrend:  "stable",
		IncomeTrend:  "stable",
		ExpenseTrend: "stable",
		GrowthRate:   0.0,
		Recommendations: []string{
			"Continue registrando suas transações regularmente",
			"Monitore seus gastos para identificar oportunidades de economia",
		},
	}

	// Add specific recommendations based on data
	if r.Summary.TotalExpenses > r.Summary.TotalIncome {
		r.Trends.Recommendations = append(r.Trends.Recommendations,
			"Considere reduzir despesas para melhorar seu lucro")
	}

	if r.Summary.Profit > 0 {
		r.Trends.Recommendations = append(r.Trends.Recommendations,
			"Excelente! Você está gerando lucro. Continue assim!")
	}
}

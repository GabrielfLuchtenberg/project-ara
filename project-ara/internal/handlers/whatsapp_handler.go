package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"project-ara/internal/models"
	"project-ara/internal/services"
)

type WhatsAppHandler struct {
	whatsappService    *services.WhatsAppService
	nlpService         *services.NLPService
	voiceService       *services.VoiceService
	ocrService         *services.OCRService
	transactionService *services.TransactionService
	userService        *services.UserService
}

func NewWhatsAppHandler(
	whatsappService *services.WhatsAppService,
	nlpService *services.NLPService,
	voiceService *services.VoiceService,
	ocrService *services.OCRService,
	transactionService *services.TransactionService,
	userService *services.UserService,
) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService:    whatsappService,
		nlpService:         nlpService,
		voiceService:       voiceService,
		ocrService:         ocrService,
		transactionService: transactionService,
		userService:        userService,
	}
}

func (h *WhatsAppHandler) HandleWebhook(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse the WhatsApp webhook
	webhookMessage, err := h.whatsappService.ParseWebhook(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook"})
		return
	}

	// Process each message
	for _, entry := range webhookMessage.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				if err := h.processMessage(message); err != nil {
					// Log error but don't fail the webhook
					fmt.Printf("Error processing message: %v\n", err)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *WhatsAppHandler) processMessage(message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
	Audio struct {
		ID       string `json:"id"`
		MimeType string `json:"mime_type"`
	} `json:"audio,omitempty"`
	Image struct {
		ID       string `json:"id"`
		MimeType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
	} `json:"image,omitempty"`
}) error {
	// Get or create user
	user, err := h.userService.GetOrCreateUser(message.From)
	if err != nil {
		return fmt.Errorf("failed to get/create user: %w", err)
	}

	// Check if user can create transactions
	canCreate, err := h.userService.CanUserCreateTransaction(user.ID.String())
	if err != nil {
		return fmt.Errorf("failed to check user permissions: %w", err)
	}

	if !canCreate {
		// Send subscription prompt
		subscriptionMessage := "Você atingiu o limite de 50 transações gratuitas. Para continuar usando o serviço, assine nosso plano premium por apenas R$ 9,90/mês."
		return h.whatsappService.SendMessage(message.From, subscriptionMessage)
	}

	// Process different message types
	switch message.Type {
	case "text":
		return h.processTextMessage(message.From, message.Text.Body, user)
	case "audio":
		return h.processAudioMessage(message.From, message.Audio.ID, user)
	case "image":
		return h.processImageMessage(message.From, message.Image.ID, user)
	default:
		return h.whatsappService.SendMessage(message.From, "Desculpe, não consegui processar esse tipo de mensagem. Envie texto, áudio ou uma foto de recibo.")
	}
}

func (h *WhatsAppHandler) processTextMessage(from, text string, user *models.User) error {
	ctx := context.Background()
	data, err := h.nlpService.ExtractTransaction(ctx, text)
	if err != nil {
		return h.whatsappService.SendMessage(from, "Desculpe, não consegui entender a transação. Tente novamente ou envie de outra forma.")
	}
	// Save transaction
	_, err = h.transactionService.CreateTransaction(user.ID.String(), data.Amount, data.Description, models.TransactionType(data.Type), models.TransactionSourceText)
	if err != nil {
		return h.whatsappService.SendMessage(from, "Erro ao registrar a transação. Tente novamente mais tarde.")
	}
	// Send summary (stub)
	return h.whatsappService.SendMessage(from, "Transação registrada! Valor: R$ "+fmt.Sprintf("%.2f", data.Amount)+" ("+data.Type+") - "+data.Description)
}

func (h *WhatsAppHandler) processAudioMessage(from, audioID string, user *models.User) error {
	ctx := context.Background()
	// For MVP, simulate audio URL (in production, fetch from WhatsApp API)
	audioURL := "https://example.com/audio/" + audioID + ".ogg"
	text, err := h.voiceService.TranscribeAudio(ctx, audioURL)
	if err != nil {
		return h.whatsappService.SendMessage(from, "Desculpe, não consegui transcrever o áudio. Tente novamente.")
	}
	return h.processTextMessage(from, text, user)
}

func (h *WhatsAppHandler) processImageMessage(from, imageID string, user *models.User) error {
	ctx := context.Background()
	// For MVP, simulate image URL (in production, fetch from WhatsApp API)
	imageURL := "https://example.com/image/" + imageID + ".jpg"
	data, err := h.ocrService.ExtractReceipt(ctx, imageURL)
	if err != nil {
		return h.whatsappService.SendMessage(from, "Desculpe, não consegui ler o recibo. Tente novamente.")
	}
	_, err = h.transactionService.CreateTransaction(user.ID.String(), data.Amount, data.Description, models.TransactionType(data.Type), models.TransactionSourceImage)
	if err != nil {
		return h.whatsappService.SendMessage(from, "Erro ao registrar a transação do recibo. Tente novamente mais tarde.")
	}
	return h.whatsappService.SendMessage(from, "Recibo processado! Valor: R$ "+fmt.Sprintf("%.2f", data.Amount)+" ("+data.Type+") - "+data.Description)
}

func (h *WhatsAppHandler) CreateTransaction(c *gin.Context) {
	// API endpoint for creating transactions
	c.JSON(http.StatusOK, gin.H{"message": "Transaction created"})
}

func (h *WhatsAppHandler) GetUserSummary(c *gin.Context) {
	userID := c.Param("id")

	summary, err := h.transactionService.GetFinancialSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *WhatsAppHandler) CreateSubscription(c *gin.Context) {
	// Placeholder for subscription creation
	// In Phase 3, this will integrate with payment gateways
	c.JSON(http.StatusOK, gin.H{"message": "Subscription created"})
}

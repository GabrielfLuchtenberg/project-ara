package services

import (
	"context"
	"fmt"
	"os"
)

type OCRService struct {
	openaiAPIKey string
}

func NewOCRService() *OCRService {
	return &OCRService{
		openaiAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// ExtractReceipt uses OpenAI Vision API to extract transaction data from a receipt image URL
func (s *OCRService) ExtractReceipt(ctx context.Context, imageURL string) (*TransactionData, error) {
	if s.openaiAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	// For MVP, we'll simulate the call (OpenAI Vision API is not public for images yet)
	// In production, use Google Cloud Vision or Gemini Vision API
	// Here, just return a stub for now
	return &TransactionData{
		Amount:      45.00,
		Type:        "income",
		Description: "Venda de cachorro-quente",
		Date:        "",
	}, nil
}

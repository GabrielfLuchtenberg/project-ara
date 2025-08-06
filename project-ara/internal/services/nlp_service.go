package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type TransactionData struct {
	Amount      float64
	Type        string // "income" or "expense"
	Description string
	Date        string // ISO8601 or empty for today
}

type NLPService struct {
	openaiAPIKey string
}

func NewNLPService() *NLPService {
	return &NLPService{
		openaiAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// ExtractTransaction uses OpenAI API to extract transaction data from Portuguese text
func (s *NLPService) ExtractTransaction(ctx context.Context, text string) (*TransactionData, error) {
	if s.openaiAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := buildPrompt(text)

	// Call OpenAI API (Chat Completions)
	body := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "Você é um assistente financeiro para MEIs brasileiros. Extraia informações de transações financeiras de textos em português. Responda apenas em JSON."},
			{"role": "user", "content": prompt},
		},
		"max_tokens": 100,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.openaiAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices from OpenAI")
	}

	// Parse JSON from model output
	var data TransactionData
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse model output: %w", err)
	}
	return &data, nil
}

func buildPrompt(text string) string {
	return "Extraia os seguintes campos do texto: tipo (income/expense), valor (float), descrição, data (se houver, senão vazio). Responda apenas em JSON.\nTexto: " + text
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// WhatsAppMessage represents a message from WhatsApp
type WhatsAppMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
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
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppResponse represents a response to WhatsApp
type WhatsAppResponse struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

type WhatsAppService struct {
	accessToken   string
	phoneNumberID string
	apiVersion    string
}

func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{
		accessToken:   os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		phoneNumberID: os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		apiVersion:    "v18.0",
	}
}

func (w *WhatsAppService) SendMessage(to, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", w.apiVersion, w.phoneNumberID)

	response := WhatsAppResponse{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: struct {
			Body string `json:"body"`
		}{
			Body: message,
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (w *WhatsAppService) ParseWebhook(body []byte) (*WhatsAppMessage, error) {
	var message WhatsAppMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &message, nil
}

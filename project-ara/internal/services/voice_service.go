package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type VoiceService struct {
	openaiAPIKey string
}

func NewVoiceService() *VoiceService {
	return &VoiceService{
		openaiAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// TranscribeAudio uses OpenAI Whisper API to transcribe audio from a URL
func (s *VoiceService) TranscribeAudio(ctx context.Context, audioURL string) (string, error) {
	if s.openaiAPIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	// Download audio file
	req, err := http.NewRequestWithContext(ctx, "GET", audioURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Prepare multipart form for Whisper API
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "audio.ogg")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, resp.Body); err != nil {
		return "", err
	}
	writer.WriteField("model", "whisper-1")
	writer.Close()

	// Call OpenAI Whisper API
	whisperReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", &buf)
	if err != nil {
		return "", err
	}
	whisperReq.Header.Set("Authorization", "Bearer "+s.openaiAPIKey)
	whisperReq.Header.Set("Content-Type", writer.FormDataContentType())

	whisperResp, err := http.DefaultClient.Do(whisperReq)
	if err != nil {
		return "", err
	}
	defer whisperResp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(whisperResp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Text, nil
}

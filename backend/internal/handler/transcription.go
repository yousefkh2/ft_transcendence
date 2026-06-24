package handler

import(
	"net/http"
	"strings"
	"log"
	"os"
	"bytes"
	"io"
	"encoding/json"

	"transcendence/backend/internal/model"

	"github.com/labstack/echo/v4"
)

func HandleTranscription(c echo.Context) error {
	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		return c.String(http.StatusInternalServerError, "openai api key missing\n")
	}

	file, header, err := c.Request().FormFile("audio")
	if err != nil {
		return c.String(http.StatusBadRequest, "audio file is required\n")
	}
	defer file.Close()

	log.Printf("received transcription upload: %s (%d bytes)", header.Filename, header.Size)

	return c.JSON(http.StatusOK, model.TranscriptionResponse{
		Text: "transcription not implemented yet",
	})
}

func HandleRealtimeTranscriptionSession(c echo.Context) error {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return c.String(http.StatusServiceUnavailable, "openai api key missing\n")
	}


	sessionBody := map[string]any{
		"expires_after": map[string]any{
			"anchor":  "created_at",
			"seconds": 600,
		},
		"session": map[string]any{
			"type": "transcription",
			"audio": map[string]any{
				"input": map[string]any{
					"transcription": map[string]any{
						"model":    "gpt-4o-transcribe",
						"language": "en",
					},
					"turn_detection": map[string]any{
						"type":                "server_vad",
						"threshold":           0.5,
						"prefix_padding_ms":   300,
						"silence_duration_ms": 500,
					},
					"noise_reduction": map[string]any{
						"type": "near_field",
					},
				},
			},
		},
	}
	requestBody, err := json.Marshal(sessionBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create openai request\n")
	}

	request, err := http.NewRequestWithContext(
		c.Request().Context(),
		http.MethodPost,
		"https://api.openai.com/v1/realtime/client_secrets",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create openai request\n")
	}

	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("openai realtime session request failed: %v", err)
		return c.String(http.StatusBadGateway, "openai realtime session request failed\n")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return c.String(http.StatusBadGateway, "failed to read openai respone\n")
	}

	if response.StatusCode >= http.StatusBadRequest {
		log.Printf("openai realtime session failed with status %d: %s",
			response.StatusCode, strings.TrimSpace(string(responseBody)))
		return c.String(http.StatusBadGateway, "openai realtime session failed\n")
	}

	return c.Blob(response.StatusCode, "application/json", responseBody)
}

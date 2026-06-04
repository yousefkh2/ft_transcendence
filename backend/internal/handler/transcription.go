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
)

func HandleTranscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		http.Error(w, "openai api key missing\n", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid multipart form\n", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "audio file is required\n", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("received transcription upload: %s (%d bytes)", header.Filename, header.Size)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.TranscriptionResponse{
		Text: "transcription not implemented yet",
	})
}

func HandleRealtimeTranscriptionSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		http.Error(w, "openai api key missing\n", http.StatusServiceUnavailable)
		return
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
		http.Error(w, "failed to create session request\n",
			http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		"https://api.openai.com/v1/realtime/client_secrets",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		http.Error(w, "failed to create openai request\n",
			http.StatusInternalServerError)
		return
	}

	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("openai realtime session request failed: %v", err)
		http.Error(w, "openai realtime session request failed\n",
			http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "failed to read openai response\n", http.StatusBadGateway)
		return
	}

	if response.StatusCode >= http.StatusBadRequest {
		log.Printf("openai realtime session failed with status %d: %s",
			response.StatusCode, strings.TrimSpace(string(responseBody)))
		http.Error(w, "openai realtime session failed\n", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(responseBody)
}

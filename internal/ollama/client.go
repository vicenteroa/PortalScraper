package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OLLAMA_ENDPOINT = "http://localhost:11434/api/generate"
	MODEL_NAME      = "deepseek-r1:1.5b"
	TIMEOUT         = 300 * time.Second // Aumentado a 5 minutos
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: TIMEOUT,
		},
	}
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func (c *Client) Generate(prompt string) (string, error) {
	// Limitar el prompt a 2000 caracteres
	if len(prompt) > 2100 {
		prompt = prompt[:2100] + "... [TRUNCADO]"
	}

	requestBody := GenerateRequest{
		Model:  MODEL_NAME,
		Prompt: prompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error creando solicitud: %w", err)
	}

	resp, err := c.httpClient.Post(OLLAMA_ENDPOINT, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error de conexión: %w\nEjecuta: docker-compose restart ollama", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error Ollama (%d): %s", resp.StatusCode, string(body))
	}

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return response.Response, nil
}

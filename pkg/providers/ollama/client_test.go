package ollama

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/EinStack/glide/pkg/providers/clients"

	"github.com/EinStack/glide/pkg/api/schemas"

	"github.com/EinStack/glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
)

func TestOllamaClient_ChatRequest(t *testing.T) {
	OllamaAIMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPayload, _ := io.ReadAll(r.Body)

		var data interface{}
		// Parse the JSON body
		err := json.Unmarshal(rawPayload, &data)
		if err != nil {
			t.Errorf("error decoding payload (%q): %v", string(rawPayload), err)
		}

		chatResponse, err := os.ReadFile(filepath.Clean("./testdata/chat.success.json"))
		if err != nil {
			t.Errorf("error reading ollama chat mock response: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(chatResponse)
		if err != nil {
			t.Errorf("error on sending chat response: %v", err)
		}
	})

	OllamaServer := httptest.NewServer(OllamaAIMock)
	defer OllamaServer.Close()

	ctx := context.Background()
	providerCfg := DefaultConfig()

	clientCfg := clients.DefaultClientConfig()

	providerCfg.Model = "llama2"

	client, err := NewClient(providerCfg, clientCfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.ChatRequest{Message: schemas.ChatMessage{
		Role:    "user",
		Content: "What's the biggest animal?",
	}}

	_, err = client.Chat(ctx, &request)

	// require.NoError(t, err)

	require.Error(t, err)
	require.Contains(t, err.Error(), "chat request failed")
}

func TestOllamaClient_ChatRequest_Non200Response(t *testing.T) {
	// Create a mock HTTP server that returns a non-OK status code
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer mockServer.Close()

	// Create a new client with the mock server URL
	client := &Client{
		httpClient: http.DefaultClient,
		chatURL:    mockServer.URL,
		config:     DefaultConfig(),
		telemetry:  telemetry.NewTelemetryMock(),
	}

	// Create a chat request payload
	payload := &ChatRequest{
		Messages: []ChatMessage{{Role: "human", Content: "Hello"}},
	}

	// Call the chatRequest function
	_, err := client.doChatRequest(context.Background(), payload)

	require.Error(t, err)
	require.Contains(t, err.Error(), "provider is not available")
}

func TestOllamaClient_ChatRequest_SuccessfulResponse(t *testing.T) {
	// Create a mock HTTP server that returns an OK status code and a sample response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response": "OK"}`))
	}))

	defer mockServer.Close()

	// Create a new client with the mock server URL
	client := &Client{
		httpClient: http.DefaultClient,
		chatURL:    mockServer.URL,
		config:     DefaultConfig(),
		telemetry:  telemetry.NewTelemetryMock(),
	}

	// Create a chat request payload
	payload := &ChatRequest{
		Messages: []ChatMessage{{Role: "human", Content: "Hello"}},
	}

	// Call the chatRequest function
	response, err := client.doChatRequest(context.Background(), payload)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "", response.ModelResponse.Message.Role)
}

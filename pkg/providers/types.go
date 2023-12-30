package providers

import (
	"net/http"
	"time"
)

type GatewayConfig struct {
	Gateway PoolsConfig `yaml:"gateway" validate:"required"`
}
type PoolsConfig struct {
	Pools []Pool `yaml:"pools" validate:"required"`
}

type Pool struct {
	Name      string     `yaml:"pool" validate:"required"`
	Balancing string     `yaml:"balancing" validate:"required"`
	Providers []Provider `yaml:"providers" validate:"required"`
}

type Provider struct {
	Name          string                 `yaml:"name" validate:"required"`
	Model         string                 `yaml:"model"`
	APIKey        string                 `yaml:"api_key" validate:"required"`
	TimeoutMs     int                    `yaml:"timeout_ms,omitempty"`
	DefaultParams map[string]interface{} `yaml:"default_params,omitempty"`
}

type ProviderVars struct {
	Name        string `yaml:"name"`
	ChatBaseURL string `yaml:"chatBaseURL"`
}

type RequestBody struct {
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	MessageHistory []string `json:"messageHistory"`
}

// Variables

var HTTPClient = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 2,
	},
}

type UnifiedAPIData struct {
	Model          string                 `json:"model"`
	APIKey         string                 `json:"api_key"`
	Params         map[string]interface{} `json:"params"`
	Message        map[string]string      `json:"message"`
	MessageHistory []map[string]string    `json:"messageHistory"`
}
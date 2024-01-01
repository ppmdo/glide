package routers

import (
	"fmt"

	"glide/pkg/providers/openai"
	"glide/pkg/routers/strategy"
)

type Config struct {
	LanguageRouters []LangRouterConfig `yaml:"language"`
}

type LangModel struct {
	ID        string `yaml:"id"`
	TimeoutMs *int   `yaml:"timeout_ms,omitempty"` // TODO: try to use Duration to bring more flexibility
	OpenAI    *openai.Config
	// Add other providers like
	// Cohere *cohere.Config
	// Anthropic *anthropic.Config
}

func (m *LangModel) validateOneProvider() error {
	providersConfigured := 0

	if m.OpenAI != nil {
		providersConfigured++
	}

	// check other providers here
	if providersConfigured == 0 {
		return fmt.Errorf("exactly one provider must be cofigured for model \"%v\", none is configured", m.ID)
	}

	if providersConfigured > 1 {
		return fmt.Errorf(
			"exactly one provider must be cofigured for model \"%v\", %v are configured",
			m.ID,
			providersConfigured,
		)
	}

	return nil
}

func (m *LangModel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain LangModel // to avoid recursion

	if err := unmarshal((*plain)(m)); err != nil {
		return err
	}

	return m.validateOneProvider()
}

type LangRouterConfig struct {
	ID              string                   `yaml:"id"`
	Enabled         bool                     `yaml:"enabled"`
	RoutingStrategy strategy.RoutingStrategy `yaml:"strategy"`
	Models          []LangModel              `yaml:"models"`
}

func DefaultLangRouterConfig() LangRouterConfig {
	return LangRouterConfig{
		Enabled:         true,
		RoutingStrategy: strategy.Priority,
	}
}

func (p *LangRouterConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = DefaultLangRouterConfig()

	type plain LangRouterConfig // to avoid recursion

	return unmarshal((*plain)(p))
}
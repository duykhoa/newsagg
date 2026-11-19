package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Provider       string
	OllamaEndpoint string
	OllamaModel    string
	GeminiAPIKey   string
	GeminiModel    string
	BatchSize      int
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("llm_provider", "ollama")
	viper.SetDefault("ollama_endpoint", "http://localhost:11434")
	viper.SetDefault("ollama_model", "qwen3.5:4b")
	viper.SetDefault("gemini_api_key", "")
	viper.SetDefault("gemini_model", "gemini-3.1-flash-lite")
	viper.SetDefault("batch_size", 25)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	return &Config{
		Provider:       viper.GetString("llm_provider"),
		OllamaEndpoint: viper.GetString("ollama_endpoint"),
		OllamaModel:    viper.GetString("ollama_model"),
		GeminiAPIKey:   viper.GetString("gemini_api_key"),
		GeminiModel:    viper.GetString("gemini_model"),
		BatchSize:      viper.GetInt("batch_size"),
	}, nil
}

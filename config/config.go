package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	SystemPromptFile string
	DatabaseURL      string
	EmbeddingDim     int
	EmbedderBaseURL  string
	EmbedderAPIKey   string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("OPENAI_BASE_URL"),
		APIKey:           os.Getenv("OPENAI_API_KEY"),
		Model:            os.Getenv("OPENAI_MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		EmbeddingDim:     atoiOR(os.Getenv("EMBEDDING_DIM"), 0),
		EmbedderBaseURL:  os.Getenv("ENBEDDING_BASE_URL"),
		EmbedderAPIKey:   os.Getenv("ENBEDDING_APIU_Key"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}

	if cfg.Model == "" {
		cfg.Model = "gpt-4o-mini"
	}

	if cfg.EmbeddingDim == 0 {
		cfg.EmbeddingDim = 768
	}

	if cfg.EmbedderBaseURL == "" {
		cfg.EmbedderBaseURL = cfg.BaseURL
		if cfg.EmbedderAPIKey == "" {
			cfg.EmbedderAPIKey = cfg.APIKey
		}
	}

	return cfg

}

func atoiOR(s string, fallback int) int {
	if s == "" {
		return fallback
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}

	return n
}

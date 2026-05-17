package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	SystemPromptFile string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("OPEN_BASE_URL"),
		APIKey:           os.Getenv("OPENAI_API_KEY"),
		Model:            OS.Getenv("OPENAI_MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
	}
}

package app

import (
	"context"
	"go_rag/chat"
	"go_rag/config"
	"go_rag/llm"
)

func Run(ctx context.Context, cfg config.Config) error {
	client := llm.New(cfg)
	return chat.RunREPL(ctx, client, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})

}

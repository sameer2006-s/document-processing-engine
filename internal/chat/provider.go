package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sameer2006-s/document-processing-engine/internal/config"
)

type LLMProvider interface {
	Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	ChatStream(ctx context.Context, systemPrompt, userPrompt string, onToken func(string) error) (string, error)
}

func NewProviderFromConfig(cfg config.ChatConfig) (LLMProvider, error) {
	switch cfg.Provider {
	case "openrouter":
		if cfg.OpenRouterKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY is required when CHAT_PROVIDER=openrouter")
		}
		return NewOpenRouterProvider(cfg.OpenRouterKey, cfg.Model, cfg.SiteURL, cfg.AppName), nil
	default:
		if cfg.GitHubToken == "" {
			return nil, fmt.Errorf("GITHUB_TOKEN is required when CHAT_PROVIDER=github")
		}
		return NewGitHubProvider(cfg.GitHubToken, cfg.Model), nil
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func readChatStream(body io.Reader, onToken func(string) error) (full string, err error) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk chatResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return "", fmt.Errorf("unmarshal chat response: %w", err)
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		deltaContent := chunk.Choices[0].Delta.Content
		if deltaContent == "" {
			continue
		}
		if err := onToken(deltaContent); err != nil {
			return full, fmt.Errorf("on token: %w", err)
		}
		full += deltaContent
	}
	if err := scanner.Err(); err != nil {
		return full, fmt.Errorf("scan chat response: %w", err)
	}
	return full, nil
}

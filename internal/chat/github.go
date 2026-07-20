package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const githubDefaultBaseURL = "https://models.github.ai/inference"

type GitHubProvider struct {
	token      string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewGitHubProvider(token, model string) *GitHubProvider {
	if model == "" {
		model = "openai/o4-mini"
	}
	return &GitHubProvider{
		token: token,
		model: model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: githubDefaultBaseURL,
	}
}

func (p *GitHubProvider) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if p.token == "" {
		return "", fmt.Errorf("github token is required")
	}

	payload := chatRequest{
		Model: p.model,
		Messages: []chatMessage{
			{Role: "developer", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	url := strings.TrimRight(p.baseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read chat response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := string(respBody)
		if len(snippet) > 500 {
			snippet = snippet[:500]
		}
		return "", fmt.Errorf("chat api status %d: %s", resp.StatusCode, snippet)
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("chat api returned no choices")
	}

	return parsed.Choices[0].Message.Content, nil
}

func (p *GitHubProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, onToken func(string) error) (full string, err error) {
	if p.token == "" {
		return "", fmt.Errorf("github token is required")
	}

	payload := chatRequest{
		Model: p.model,
		Messages: []chatMessage{
			{Role: "developer", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	url := strings.TrimRight(p.baseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		snippet := string(respBody)
		if len(snippet) > 500 {
			snippet = snippet[:500]
		}
		return "", fmt.Errorf("chat api status %d: %s", resp.StatusCode, snippet)
	}

	return readChatStream(resp.Body, onToken)
}

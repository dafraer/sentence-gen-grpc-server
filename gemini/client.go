package gemini

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Client struct {
	client *genai.Client
	logger *zap.SugaredLogger
}

// New creates new gemini client
func New(ctx context.Context, token string, logger *zap.SugaredLogger) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, err
	}
	return &Client{client: client, logger: logger}, nil
}

// Close closes the client
func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

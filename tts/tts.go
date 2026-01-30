package tts

import (
	"context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"go.uber.org/zap"
)

type Client struct {
	tts    *texttospeech.Client
	logger *zap.SugaredLogger
}

// New creates new tts client
func New(ctx context.Context, logger *zap.SugaredLogger) (*Client, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		tts:    client,
		logger: logger,
	}, nil
}

// Close closes tts client
func (c *Client) Close() error {
	if err := c.tts.Close(); err != nil {
		return err
	}
	return nil
}

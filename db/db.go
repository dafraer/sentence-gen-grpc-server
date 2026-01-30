package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type Store struct {
	db     *firestore.Client
	logger *zap.SugaredLogger
}

// New creates new firestore instance
func New(ctx context.Context, logger *zap.SugaredLogger, projectID string) (*Store, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Store{db: client, logger: logger}, nil
}

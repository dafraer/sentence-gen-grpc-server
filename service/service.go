package service

import (
	"github.com/dafraer/sentence-gen-grpc-server/db"
	"github.com/dafraer/sentence-gen-grpc-server/gemini"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
	"go.uber.org/zap"
)

type Service struct {
	ttsClient    *tts.Client
	geminiClient *gemini.Client
	logger       *zap.SugaredLogger
	store        *db.Store
}

func New(ttsClient *tts.Client, geminiClient *gemini.Client, logger *zap.SugaredLogger, store *db.Store) *Service {
	return &Service{
		ttsClient:    ttsClient,
		geminiClient: geminiClient,
		logger:       logger,
		store:        store,
	}
}

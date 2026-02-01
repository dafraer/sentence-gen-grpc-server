package service

import (
	"github.com/dafraer/sentence-gen-grpc-server/db"
	"github.com/dafraer/sentence-gen-grpc-server/gemini"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
	"go.uber.org/zap"
)

var (
	ErrEmpty
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

func (s *Service) GenerateSentence(req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	if req.Word == "" || req.WordLanguage == "" || req.TranslationLanguage == "" {

	}
	if req.IncludeAudio {

	}
	return nil, nil
}

func (s *Service) Translate(req *TranslateRequest) (*TranslateResponse, error) {
	return nil, nil
}

func (s *Service) GenerateDefinition(req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	return nil, nil
}

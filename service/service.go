package service

import (
	"context"
	"errors"

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

func (s *Service) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	sentences, err := s.geminiClient.GenerateSentence(ctx, &gemini.SentenceGenerationRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
	})
	if err != nil {
		return nil, err
	}

	resp := &GenerateSentenceResponse{
		OriginalSentence:   sentences.OriginalSentence,
		TranslatedSentence: sentences.TranslatedSentence,
	}

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		audio, err := s.ttsClient.Generate(ctx, sentences.OriginalSentence, req.WordLanguage, gender, tts.Chirp3HD)
		if err != nil {
			return nil, err
		}
		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	translation, err := s.geminiClient.Translate(ctx, &gemini.TranslationRequest{
		Word:            req.Word,
		FromLanguage:    req.FromLanguage,
		ToLanguage:      req.ToLanguage,
		TranslationHint: req.TranslationHint,
	})
	if err != nil {
		return nil, err
	}

	resp := &TranslateResponse{
		Translation: translation.Translation,
	}

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		audio, err := s.ttsClient.Generate(ctx, req.Word, req.FromLanguage, gender, tts.Chirp3HD)
		if err != nil {
			return nil, err
		}
		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	definition, err := s.geminiClient.GenerateDefinition(ctx, &gemini.DefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
	})
	if err != nil {
		return nil, err
	}
	resp := &GenerateDefinitionResponse{
		Definition: definition.Definition,
	}

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		audio, err := s.ttsClient.Generate(ctx, req.Word, req.Language, gender, tts.Chirp3HD)
		if err != nil && !errors.Is(err, tts.ErrNoSuchVoice) {
			return nil, err
		}
		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		return nil, err
	}

	return resp, nil
}

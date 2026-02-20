package service

import (
	"context"
	"errors"

	"github.com/dafraer/sentence-gen-grpc-server/config"
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
	config       *config.Config
}

func New(ttsClient *tts.Client, geminiClient *gemini.Client, logger *zap.SugaredLogger, store *db.Store, cfg *config.Config) *Service {
	return &Service{
		ttsClient:    ttsClient,
		geminiClient: geminiClient,
		logger:       logger,
		store:        store,
		config:       cfg,
	}
}

func (s *Service) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	s.logger.Infow("generate sentence request received", "word", req.Word, "word_language", req.WordLanguage, "translation_language", req.TranslationLanguage, "include_audio", req.IncludeAudio)

	if err := req.validate(); err != nil {
		s.logger.Errorw("generate sentence request validation failed", "error", err)
		return nil, err
	}

	sentences, tokenCnt, err := s.geminiClient.GenerateSentence(ctx, &gemini.SentenceGenerationRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
	})
	if err != nil {
		s.logger.Errorw("generate sentence via gemini failed", "error", err)
		return nil, err
	}
	s.logger.Debugw("generate sentence via gemini succeeded", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)

	resp := &GenerateSentenceResponse{
		OriginalSentence:   sentences.OriginalSentence,
		TranslatedSentence: sentences.TranslatedSentence,
	}

	if err := s.AddSpending(ctx, &AddDailySpendingParams{
		GeminiInputTokens:  tokenCnt.InputTokens,
		GeminiOutputTokens: tokenCnt.OutputTokens,
	}); err != nil {
		s.logger.Errorw("failed to add gemini spending for sentence generation", "error", err)
		return nil, err
	}
	s.logger.Debugw("added gemini spending for sentence generation", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		s.logger.Debugw("generating sentence audio", "language", req.WordLanguage, "gender", gender, "model", tts.Chirp3HD)
		audio, err := s.ttsClient.Generate(ctx, sentences.OriginalSentence, req.WordLanguage, gender, tts.Chirp3HD) //TODO: Should be variable in the future
		if err != nil && !errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Errorw("sentence audio generation failed", "error", err)
			return nil, err
		}
		if errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Debugw("sentence audio generation skipped due to missing voice", "language", req.WordLanguage, "gender", gender, "model", tts.Chirp3HD)
		}

		if !errors.Is(err, tts.ErrNoSuchVoice) {
			if err := s.AddSpending(ctx, &AddDailySpendingParams{
				Characters: int64(len([]rune(sentences.OriginalSentence))),
				TTSModel:   tts.Chirp3HD, //TODO: Should be variable in the future
			}); err != nil {
				s.logger.Errorw("failed to add tts spending for sentence generation", "error", err)
				return nil, err
			}
			s.logger.Debugw("added tts spending for sentence generation", "characters", int64(len([]rune(sentences.OriginalSentence))), "model", tts.Chirp3HD)
		}
		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		s.logger.Errorw("generate sentence response validation failed", "error", err)
		return nil, err
	}
	s.logger.Infow("generate sentence request completed", "has_audio", len(resp.Audio) > 0)

	return resp, nil
}

func (s *Service) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	s.logger.Infow("translate request received", "word", req.Word, "from_language", req.FromLanguage, "to_language", req.ToLanguage, "include_audio", req.IncludeAudio)

	if err := req.validate(); err != nil {
		s.logger.Errorw("translate request validation failed", "error", err)
		return nil, err
	}

	translation, tokenCnt, err := s.geminiClient.Translate(ctx, &gemini.TranslationRequest{
		Word:            req.Word,
		FromLanguage:    req.FromLanguage,
		ToLanguage:      req.ToLanguage,
		TranslationHint: req.TranslationHint,
	})
	if err != nil {
		s.logger.Errorw("translate via gemini failed", "error", err)
		return nil, err
	}
	s.logger.Debugw("translate via gemini succeeded", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)

	resp := &TranslateResponse{
		Translation: translation.Translation,
	}

	if err := s.AddSpending(ctx, &AddDailySpendingParams{
		GeminiInputTokens:  tokenCnt.InputTokens,
		GeminiOutputTokens: tokenCnt.OutputTokens,
	}); err != nil {
		s.logger.Errorw("failed to add gemini spending for translation", "error", err)
		return nil, err
	}
	s.logger.Debugw("added gemini spending for translation", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		s.logger.Debugw("generating translation audio", "language", req.FromLanguage, "gender", gender, "model", tts.Chirp3HD)
		audio, err := s.ttsClient.Generate(ctx, req.Word, req.FromLanguage, gender, tts.Chirp3HD)
		if err != nil && !errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Errorw("translation audio generation failed", "error", err)
			return nil, err
		}
		if errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Debugw("translation audio generation skipped due to missing voice", "language", req.FromLanguage, "gender", gender, "model", tts.Chirp3HD)
		}

		if !errors.Is(err, tts.ErrNoSuchVoice) {
			if err := s.AddSpending(ctx, &AddDailySpendingParams{
				Characters: int64(len([]rune(req.Word))),
				TTSModel:   tts.Chirp3HD, //TODO: Should be variable in the future
			}); err != nil {
				s.logger.Errorw("failed to add tts spending for translation", "error", err)
				return nil, err
			}
			s.logger.Debugw("added tts spending for translation", "characters", int64(len([]rune(req.Word))), "model", tts.Chirp3HD)
		}

		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		s.logger.Errorw("translate response validation failed", "error", err)
		return nil, err
	}
	s.logger.Infow("translate request completed", "has_audio", len(resp.Audio) > 0)

	return resp, nil
}

func (s *Service) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	s.logger.Infow("generate definition request received", "word", req.Word, "language", req.Language, "include_audio", req.IncludeAudio)

	if err := req.validate(); err != nil {
		s.logger.Errorw("generate definition request validation failed", "error", err)
		return nil, err
	}
	definition, tokenCnt, err := s.geminiClient.GenerateDefinition(ctx, &gemini.DefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
	})
	if err != nil {
		s.logger.Errorw("generate definition via gemini failed", "error", err)
		return nil, err
	}
	s.logger.Debugw("generate definition via gemini succeeded", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)
	resp := &GenerateDefinitionResponse{
		Definition: definition.Definition,
	}

	if err := s.AddSpending(ctx, &AddDailySpendingParams{
		GeminiInputTokens:  tokenCnt.InputTokens,
		GeminiOutputTokens: tokenCnt.OutputTokens,
	}); err != nil {
		s.logger.Errorw("failed to add gemini spending for definition generation", "error", err)
		return nil, err
	}
	s.logger.Debugw("added gemini spending for definition generation", "input_tokens", tokenCnt.InputTokens, "output_tokens", tokenCnt.OutputTokens)

	if req.IncludeAudio {
		gender := tts.Female
		if req.VoiceGender == Male {
			gender = tts.Male
		}
		s.logger.Debugw("generating definition audio", "language", req.Language, "gender", gender, "model", tts.Chirp3HD)
		audio, err := s.ttsClient.Generate(ctx, req.Word, req.Language, gender, tts.Chirp3HD)
		if err != nil && !errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Errorw("definition audio generation failed", "error", err)
			return nil, err
		}
		if errors.Is(err, tts.ErrNoSuchVoice) {
			s.logger.Debugw("definition audio generation skipped due to missing voice", "language", req.Language, "gender", gender, "model", tts.Chirp3HD)
		}

		if !errors.Is(err, tts.ErrNoSuchVoice) {
			if err := s.AddSpending(ctx, &AddDailySpendingParams{
				Characters: int64(len([]rune(req.Word))),
				TTSModel:   tts.Chirp3HD, //TODO: Should be variable in the future
			}); err != nil {
				s.logger.Errorw("failed to add tts spending for definition generation", "error", err)
				return nil, err
			}
			s.logger.Debugw("added tts spending for definition generation", "characters", int64(len([]rune(req.Word))), "model", tts.Chirp3HD)
		}

		resp.Audio = audio
	}

	if err := resp.validate(); err != nil {
		s.logger.Errorw("generate definition response validation failed", "error", err)
		return nil, err
	}
	s.logger.Infow("generate definition request completed", "has_audio", len(resp.Audio) > 0)

	return resp, nil
}

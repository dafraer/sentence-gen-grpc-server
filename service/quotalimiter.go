package service

import (
	"context"
	"errors"

	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"github.com/dafraer/sentence-gen-grpc-server/db"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
)

const (
	chirp3HDVoicePerCharacterPrice = currency.MicroUSD(30)
	standardVoicePerCharacterPrice = currency.MicroUSD(4)
)

func (s *Service) DailyQuotaExceeded(ctx context.Context) (bool, error) {
	s.logger.Debugw("checking daily quota")
	spending, err := s.store.GetDailySpending(ctx)
	if err != nil {
		s.logger.Errorw("failed to get daily spending", "error", err)
		return false, err
	}
	if spending.Amount >= s.config.DailyQuota {
		s.logger.Infow("daily quota exceeded", "amount", spending.Amount, "quota", s.config.DailyQuota)
		return true, nil
	}
	s.logger.Debugw("daily quota check passed", "amount", spending.Amount, "quota", s.config.DailyQuota)
	return false, nil
}

func (s *Service) AddSpending(ctx context.Context, params *AddDailySpendingParams) error {
	if params == nil {
		s.logger.Errorw("add spending failed: nil params", "error", errors.New("params cannot be nil"))
		return errors.New("params cannot be nil")
	}

	sp := db.Spending{}
	if params.TTSModel == tts.Chirp3HD {
		sp.Amount += currency.MicroUSD(params.Characters) * chirp3HDVoicePerCharacterPrice
		sp.Chirp3HDCharacters += params.Characters
	}

	if params.TTSModel == tts.Standard {
		sp.Amount += currency.MicroUSD(params.Characters) * standardVoicePerCharacterPrice
		sp.StandardVoiceCharacters += params.Characters
	}

	sp.Amount += s.config.GeminiInputPrice * currency.MicroUSD(params.GeminiInputTokens)
	sp.GeminiInputTokens = params.GeminiInputTokens

	sp.Amount += s.config.GeminiOutputPrice * currency.MicroUSD(params.GeminiOutputTokens)
	sp.GeminiOutputTokens = params.GeminiOutputTokens

	if err := s.store.AddDailySpending(ctx, &db.Spending{
		Amount:                  sp.Amount,
		Chirp3HDCharacters:      sp.Chirp3HDCharacters,
		StandardVoiceCharacters: sp.StandardVoiceCharacters,
		GeminiInputTokens:       sp.GeminiInputTokens,
		GeminiOutputTokens:      sp.GeminiOutputTokens,
	}); err != nil {
		s.logger.Errorw("failed to persist spending", "error", err)
		return err
	}
	s.logger.Debugw("spending persisted", "amount", sp.Amount, "chirp3hd_characters", sp.Chirp3HDCharacters, "standard_characters", sp.StandardVoiceCharacters, "gemini_input_tokens", sp.GeminiInputTokens, "gemini_output_tokens", sp.GeminiOutputTokens)

	return nil
}

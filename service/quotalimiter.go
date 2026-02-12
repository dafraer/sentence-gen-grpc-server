package service

import (
	"context"
	"errors"

	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
)

const (
	Chirp3HDVoicePerCharacterPrice = currency.MicroUSD(30)
	StandardVoicePerCharacterPrice = currency.MicroUSD(4)
)

func (s *Service) DailyQuotaExceeded(ctx context.Context) (bool, error) {
	spending, err := s.store.GetDailySpending(ctx)
	if err != nil {
		return false, err
	}
	if spending.Amount > s.config.DailyQuota {
		return true, nil
	}
	return false, nil
}

func (s *Service) AddSpending(ctx context.Context, params *UpdateDailySpendingParams) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	spending, err := s.store.GetDailySpending(ctx)
	if err != nil {
		return err
	}

	if params.TTSModel == tts.Chirp3HD {
		spending.Amount += currency.MicroUSD(params.Characters) * Chirp3HDVoicePerCharacterPrice
		spending.StandardVoiceCharacters += params.Characters
	}

	if params.TTSModel == tts.Standard {
		spending.Amount += currency.MicroUSD(params.Characters) * StandardVoicePerCharacterPrice
		spending.StandardVoiceCharacters += params.Characters
	}

	spending.Amount += s.config.GeminiInputPrice * currency.MicroUSD(params.GeminiInputTokens)
	spending.GeminiInputTokens += params.GeminiInputTokens

	spending.Amount += s.config.GeminiOutputPrice * currency.MicroUSD(params.GeminiOutputTokens)
	spending.GeminiOutputTokens += params.GeminiOutputTokens

	if err := s.store.UpdateDailySpending(ctx, spending); err != nil {
		return err
	}

	if err := s.store.UpdateTotalSpending(spending); err != nil {
		return err
	}
	return nil
}

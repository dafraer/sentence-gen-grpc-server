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
	standardVoiceFreeLimit         = 4_000_000
	chirp3HDVoiceFreeLimit         = 1_000_000
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

	dailySpending, err := s.store.GetDailySpending(ctx)
	if err != nil {
		return err
	}

	totalSpending, err := s.store.GetTotalSpending(ctx)
	if err != nil {
		return err
	}

	sp := db.Spending{}
	//TODO:Fix incorrect calculation of free tier
	if params.TTSModel == tts.Chirp3HD {
		if totalSpending.Chirp3HDCharacters > chirp3HDVoiceFreeLimit {
			sp.Amount += currency.MicroUSD(params.Characters) * chirp3HDVoicePerCharacterPrice
		}
		sp.Chirp3HDCharacters += params.Characters
	}

	if params.TTSModel == tts.Standard {
		if totalSpending.StandardVoiceCharacters > standardVoiceFreeLimit {
			sp.Amount += currency.MicroUSD(params.Characters) * standardVoicePerCharacterPrice
		}
		sp.StandardVoiceCharacters += params.Characters
	}

	sp.Amount += s.config.GeminiInputPrice * currency.MicroUSD(params.GeminiInputTokens)
	sp.GeminiInputTokens += params.GeminiInputTokens

	sp.Amount += s.config.GeminiOutputPrice * currency.MicroUSD(params.GeminiOutputTokens)
	sp.GeminiOutputTokens += params.GeminiOutputTokens

	if err := s.store.UpdateDailySpending(ctx, &db.Spending{
		Amount:                  dailySpending.Amount + sp.Amount,
		Chirp3HDCharacters:      dailySpending.Chirp3HDCharacters + sp.Chirp3HDCharacters,
		StandardVoiceCharacters: dailySpending.StandardVoiceCharacters + sp.StandardVoiceCharacters,
		GeminiInputTokens:       dailySpending.GeminiInputTokens + sp.GeminiInputTokens,
		GeminiOutputTokens:      dailySpending.GeminiOutputTokens + sp.GeminiOutputTokens,
	}); err != nil {
		return err
	}

	if err := s.store.UpdateTotalSpending(ctx, &db.Spending{
		Amount:                  totalSpending.Amount + sp.Amount,
		Chirp3HDCharacters:      totalSpending.Chirp3HDCharacters + sp.Chirp3HDCharacters,
		StandardVoiceCharacters: totalSpending.StandardVoiceCharacters + sp.StandardVoiceCharacters,
		GeminiInputTokens:       totalSpending.GeminiInputTokens + sp.GeminiInputTokens,
		GeminiOutputTokens:      totalSpending.GeminiOutputTokens + sp.GeminiOutputTokens,
	}); err != nil {
		return err
	}
	return nil
}

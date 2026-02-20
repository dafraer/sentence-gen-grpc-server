package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func initDB() (*Store, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	store, err := New(context.Background(), sugar, "enhanced-rarity-437111-d9")
	if err != nil {
		return nil, err
	}
	return store, nil
}

func TestStore_GetDailySpending(t *testing.T) {
	s, err := initDB()
	assert.NoError(t, err)
	res, err := s.GetDailySpending(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NoError(t, s.Close())
}

func TestStore_AddDailySpending(t *testing.T) {
	s, err := initDB()
	assert.NoError(t, err)
	assert.NoError(t, s.AddDailySpending(context.Background(), &Spending{
		Amount:                  10,
		Chirp3HDCharacters:      10,
		StandardVoiceCharacters: 10,
		GeminiInputTokens:       10,
		GeminiOutputTokens:      10,
	}))
	assert.NoError(t, s.Close())
}

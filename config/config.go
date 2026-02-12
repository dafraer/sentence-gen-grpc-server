package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"github.com/joho/godotenv"
)

type Config struct {
	DailyQuota        currency.MicroUSD
	ProjectID         string
	Address           string
	GeminiModel       string
	GeminiInputPrice  currency.MicroUSD
	GeminiOutputPrice currency.MicroUSD
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	quota, err := strconv.Atoi(os.Getenv("DAILY_QUOTA"))
	if err != nil {
		return nil, err
	}

	inputPrice, err := strconv.Atoi(os.Getenv("GEMINI_INPUT_PRICE"))
	if err != nil {
		return nil, err
	}

	outputPrice, err := strconv.Atoi(os.Getenv("GEMINI_OUTPUT_PRICE"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		GeminiInputPrice:  currency.MicroUSD(inputPrice),
		GeminiOutputPrice: currency.MicroUSD(outputPrice),
		DailyQuota:        currency.MicroUSD(quota),
		ProjectID:         os.Getenv("PROJECT_ID"),
		Address:           os.Getenv("ADDRESS"),
		GeminiModel:       os.Getenv("GEMINI_MODEL"),
	}
	if cfg.DailyQuota == 0 || cfg.ProjectID == "" || cfg.Address == "" || cfg.GeminiModel == "" {
		return nil, errors.New("invalid configuration")
	}
	return cfg, nil
}

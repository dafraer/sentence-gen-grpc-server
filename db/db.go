package db

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"go.uber.org/zap"
)

const (
	collectionSpending    = "spending"
	documentTotal         = "total"
	amountKey             = "amount_micro_usd"
	chirp3HDCharsKey      = "chirp3hd_characters"
	standardVoiceCharsKey = "standard_voice_characters"
	geminiInputTokensKey  = "gemini_input_tokens"
	geminiOutputTokensKey = "gemini_output_tokens"
)

type Store struct {
	db     *firestore.Client
	logger *zap.SugaredLogger
}

type Spending struct {
	Amount                  currency.MicroUSD `firestore:"amount_micro_usd"`
	Chirp3HDCharacters      int               `firestore:"chirp3hd_characters"`
	StandardVoiceCharacters int               `firestore:"standard_voice_characters"`
	GeminiInputTokens       int               `firestore:"gemini_input_tokens"`
	GeminiOutputTokens      int               `firestore:"gemini_output_tokens"`
}

// New creates new firestore instance
func New(ctx context.Context, logger *zap.SugaredLogger, projectID string) (*Store, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return &Store{db: client, logger: logger}, nil
}

func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetDailySpending(ctx context.Context) (*Spending, error) {
	day := time.Now().In(time.UTC).Format("2006-01-02")
	docSnap, err := s.db.Collection(collectionSpending).Doc(day).Get(ctx)
	if err != nil {
		return nil, err
	}

	var sp Spending
	if err := docSnap.DataTo(&sp); err != nil {
		return nil, err
	}
	return &sp, nil
}

func (s *Store) GetTotalSpending(ctx context.Context) (*Spending, error) {
	docSnap, err := s.db.Collection(collectionSpending).Doc(documentTotal).Get(ctx)
	if err != nil {
		return nil, err
	}

	var sp Spending
	if err := docSnap.DataTo(&sp); err != nil {
		return nil, err
	}
	return &sp, nil
}

// AddDailySpending updates all the fields in the daily spending doc
func (s *Store) AddDailySpending(ctx context.Context, params *Spending) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	day := time.Now().In(time.UTC).Format("2006-01-02")
	doc := s.db.Collection(collectionSpending).Doc(day)

	_, err := doc.Set(ctx, map[string]interface{}{
		amountKey:             firestore.Increment(int(params.Amount)),
		chirp3HDCharsKey:      firestore.Increment(params.Chirp3HDCharacters),
		standardVoiceCharsKey: firestore.Increment(params.StandardVoiceCharacters),
		geminiInputTokensKey:  firestore.Increment(params.GeminiInputTokens),
		geminiOutputTokensKey: firestore.Increment(params.GeminiOutputTokens),
	}, firestore.MergeAll)

	return err
}

// AddTotalSpending updates all the fields in the daily spending doc
func (s *Store) AddTotalSpending(ctx context.Context, params *Spending) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	doc := s.db.Collection(collectionSpending).Doc(documentTotal)

	_, err := doc.Set(ctx, map[string]interface{}{
		amountKey:             firestore.Increment(int(params.Amount)),
		chirp3HDCharsKey:      firestore.Increment(params.Chirp3HDCharacters),
		standardVoiceCharsKey: firestore.Increment(params.StandardVoiceCharacters),
		geminiInputTokensKey:  firestore.Increment(params.GeminiInputTokens),
		geminiOutputTokensKey: firestore.Increment(params.GeminiOutputTokens),
	}, firestore.MergeAll)

	return err
}

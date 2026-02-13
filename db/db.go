package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"go.uber.org/zap"
)

const (
	collectionSpending    = "spending"
	documentDaily         = "daily"
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
		log.Fatalln(err)
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
	docSnap, err := s.db.Collection(collectionSpending).Doc(documentDaily).Get(ctx)
	if err != nil {
		return nil, err
	}

	var sp Spending
	if err := docSnap.DataTo(&sp); err != nil {
		return nil, err
	}
	return &sp, nil
}

// UpdateDailySpending updates all the fields in the daily spending doc
func (s *Store) UpdateDailySpending(ctx context.Context, params *Spending) error {
	_, err := s.db.
		Collection(collectionSpending).
		Doc(documentDaily).
		Set(ctx, params)
	return err
}

// UpdateTotalSpending updates all the fields in the total spending doc
func (s *Store) UpdateTotalSpending(ctx context.Context, params *Spending) error {
	_, err := s.db.
		Collection(collectionSpending).
		Doc(documentTotal).
		Set(ctx, params)
	return err
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

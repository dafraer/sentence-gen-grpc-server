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
	collectionName = "spending"
	documentDaily  = "daily"
	documentTotal  = "total"
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
	dsnap, err := s.db.Collection(collectionName).Doc(documentDaily).Get(ctx)
	if err != nil {
		return nil, err
	}
	m := dsnap.Data()
	_ = m
	return &Spending{}, nil
}

func (s *Store) UpdateDailySpending(ctx context.Context, params *Spending) error {
	doc := s.db.Collection("spending").Doc("daily")
	_, err := doc.Set(ctx, params, firestore.MergeAll)
	return err
}

func (s *Store) DeleteDailySpending() error {
	return nil
}

func (s *Store) UpdateTotalSpending(params *Spending) error {
	return nil
}

func (s *Store) GetTotalSpending() (*Spending, error) {
	return nil, nil
}

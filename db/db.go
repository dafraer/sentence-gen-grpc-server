package db

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/dafraer/sentence-gen-grpc-server/currency"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionSpending    = "spending"
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
	Chirp3HDCharacters      int64             `firestore:"chirp3hd_characters"`
	StandardVoiceCharacters int64             `firestore:"standard_voice_characters"`
	GeminiInputTokens       int64             `firestore:"gemini_input_tokens"`
	GeminiOutputTokens      int64             `firestore:"gemini_output_tokens"`
}

// New creates new firestore instance
func New(ctx context.Context, logger *zap.SugaredLogger, projectID string) (*Store, error) {
	logger.Infow("initializing firestore client", "project_id", projectID)
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		logger.Errorw("failed to initialize firebase app", "error", err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Errorw("failed to initialize firestore client", "error", err)
		return nil, err
	}
	logger.Infow("firestore client initialized")
	return &Store{db: client, logger: logger}, nil
}

func (s *Store) Close() error {
	s.logger.Infow("closing firestore client")
	if err := s.db.Close(); err != nil {
		s.logger.Errorw("failed to close firestore client", "error", err)
		return err
	}
	s.logger.Debugw("firestore client closed")
	return nil
}

func (s *Store) GetDailySpending(ctx context.Context) (*Spending, error) {
	day := time.Now().In(time.UTC).Format("2006-01-02")

	s.logger.Debugw("fetching daily spending", "day", day)

	docSnap, err := s.db.Collection(collectionSpending).Doc(day).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.logger.Debugw("daily spending not found, returning zero spending", "day", day)
			return &Spending{}, nil
		}
		s.logger.Errorw("failed to fetch daily spending", "error", err)
		return nil, err
	}

	var sp Spending
	if err := docSnap.DataTo(&sp); err != nil {
		s.logger.Errorw("failed to decode daily spending", "error", err)
		return nil, err
	}
	s.logger.Debugw("fetched daily spending", "day", day, "amount", sp.Amount)
	return &sp, nil
}

// AddDailySpending updates all the fields in the daily spending doc
func (s *Store) AddDailySpending(ctx context.Context, params *Spending) error {
	if params == nil {
		s.logger.Errorw("failed to add daily spending: nil params", "error", errors.New("params cannot be nil"))
		return errors.New("params cannot be nil")
	}

	day := time.Now().In(time.UTC).Format("2006-01-02")
	doc := s.db.Collection(collectionSpending).Doc(day)
	s.logger.Debugw("adding daily spending", "day", day, "amount", params.Amount, "chirp3hd_characters", params.Chirp3HDCharacters, "standard_characters", params.StandardVoiceCharacters, "gemini_input_tokens", params.GeminiInputTokens, "gemini_output_tokens", params.GeminiOutputTokens)

	_, err := doc.Set(ctx, map[string]interface{}{
		amountKey:             firestore.Increment(int64(params.Amount)),
		chirp3HDCharsKey:      firestore.Increment(params.Chirp3HDCharacters),
		standardVoiceCharsKey: firestore.Increment(params.StandardVoiceCharacters),
		geminiInputTokensKey:  firestore.Increment(params.GeminiInputTokens),
		geminiOutputTokensKey: firestore.Increment(params.GeminiOutputTokens),
	}, firestore.MergeAll)

	if err != nil {
		s.logger.Errorw("failed to add daily spending", "error", err)
		return err
	}
	s.logger.Debugw("daily spending added", "day", day)

	return nil
}

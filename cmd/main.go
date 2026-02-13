package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/dafraer/sentence-gen-grpc-server/config"
	"github.com/dafraer/sentence-gen-grpc-server/db"
	"github.com/dafraer/sentence-gen-grpc-server/gemini"
	"github.com/dafraer/sentence-gen-grpc-server/server"
	"github.com/dafraer/sentence-gen-grpc-server/service"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	//Declare context that is marked Done when os.Interrupt is called
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//Create firestore client
	store, err := db.New(ctx, sugar, cfg.ProjectID)
	if err != nil {
		panic(err)
	}
	defer func(store *db.Store) {
		if err := store.Close(); err != nil {
			panic(err)
		}
	}(store)
	if err := store.UpdateDailySpending(ctx, &db.Spending{StandardVoiceCharacters: 10}); err != nil {
		panic(err)
	}
	spending, err := store.GetDailySpending(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(*spending)
	return

	//Create gemini client
	geminiClient, err := gemini.New(ctx, sugar, cfg.GeminiModel)
	if err != nil {
		panic(err)
	}

	//Create tts client
	ttsClient, err := tts.New(ctx, sugar)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ttsClient.Close(); err != nil {
			panic(err)
		}
	}()

	//Create new service
	srvc := service.New(ttsClient, geminiClient, sugar, store, cfg)

	//Create new grpc server
	srv := server.NewServer(srvc, sugar)

	//Run the server
	if err := srv.Run(ctx, cfg.Address); err != nil {
		panic(err)
	}
}

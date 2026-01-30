package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/dafraer/sentence-gen-grpc-server/db"
	"github.com/dafraer/sentence-gen-grpc-server/gemini"
	"github.com/dafraer/sentence-gen-grpc-server/server"
	"github.com/dafraer/sentence-gen-grpc-server/service"
	"github.com/dafraer/sentence-gen-grpc-server/tts"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	projectID := os.Getenv("PROJECT_ID")
	port := os.Getenv("PORT")
	if port == "" || geminiAPIKey == "" || projectID == "" {
		panic("Missing GEMINI_API_KEY, PROJECT_ID or PORT")
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
	store, err := db.New(ctx, sugar, projectID)
	if err != nil {
		panic(err)
	}

	//Create gemini client
	geminiClient, err := gemini.New(ctx, geminiAPIKey, sugar)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := geminiClient.Close(); err != nil {
			panic(err)
		}
	}()

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
	srvc := service.New(ttsClient, geminiClient, sugar, store)

	//Create new grpc server
	srv := server.NewServer(srvc, sugar)

	//Run the server
	if err := srv.Run(ctx); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Southclaws/scribble.rs/src/app"
)

func main() {
	zap.L().Info("service initialising")
	a, err := app.Initialise(context.Background())
	if err != nil {
		zap.L().Fatal("service initialisation failed", zap.Error(err))
	}
	zap.L().Info("service initialised")
	if err := a.Start(); err != nil {
		zap.L().Fatal("service terminated", zap.Error(err))
	}
	zap.L().Info("service terminated gracefully")
}

func init() {
	godotenv.Load("../.env", ".env")

	var (
		debug  = os.Getenv("DEBUG")
		config zap.Config
		err    error
	)

	config = zap.NewProductionConfig()

	if debug != "0" && debug != "" {
		config.Level = zap.NewAtomicLevel()
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config.Level.SetLevel(zap.InfoLevel)
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}

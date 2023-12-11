package main

import (
	"context"
	"log"
	"net/http"

	"github.com/tkhq/oauth/internal/configs"
	"github.com/tkhq/oauth/internal/handlers"
	"github.com/tkhq/oauth/internal/helpers/google"
	"github.com/tkhq/oauth/internal/logger"

	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper across the application
	err := configs.InitializeViper()
	if err != nil {
		log.Fatalf("error while initializing configuration (viper): %s", err.Error())
	}

	// Initialize Logger across the application
	logger, err := logger.CreateZapLogger()
	if err != nil {
		log.Fatalf("error while initializing logger: %s", err.Error())
	}

	ctx := context.Background()

	oidc, err := google.NewGoogleOIDC(ctx)
	if err != nil {
		log.Fatalf("error while initializing Google OIDC: %s", err.Error())
	}

	// Routes for the application
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/oidc/google", handlers.LoginHandler(logger, oidc))
	http.HandleFunc("/oidc/google/callback", handlers.CallbackHandler(ctx, oidc, logger))

	logger.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))
}

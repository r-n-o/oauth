package google

import (
	"context"
	"errors"

	goOidc "github.com/coreos/go-oidc"
	"go.uber.org/zap"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOIDC struct {
	ClientId    string
	RedirectURL string
	provider    *goOidc.Provider
	OauthConfig *oauth2.Config
	OauthState  string

	// must keep private
	clientSecret string
}

type Claims struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}

func NewGoogleOIDC(ctx context.Context) (*GoogleOIDC, error) {
	provider, err := goOidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &GoogleOIDC{
		ClientId:     viper.GetString("google.clientID"),
		RedirectURL:  "",
		clientSecret: viper.GetString("google.clientSecret"),
		provider:     provider,
		OauthConfig: &oauth2.Config{
			ClientID:     viper.GetString("google.clientID"),
			ClientSecret: viper.GetString("google.clientSecret"),
			RedirectURL:  "http://localhost:9090/oidc/google/callback",
			Scopes:       []string{goOidc.ScopeOpenID, "profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}, nil
}

// Does an oauth exchange and get the OIDC token (string)
func (oidc *GoogleOIDC) GetToken(ctx context.Context, code string, logger *zap.Logger) (string, error) {
	token, err := oidc.OauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	logger.Info("TOKEN>> AccessToken>> " + token.AccessToken)
	logger.Info("TOKEN>> Expiration Time>> " + token.Expiry.String())
	logger.Info("TOKEN>> RefreshToken>> " + token.RefreshToken)

	rawIDToken, ok := token.Extra("id_token").(string)
	logger.Info("RAW ID TOKEN>> " + rawIDToken)
	if !ok {
		return "", errors.New("no ID token found")
	}

	return rawIDToken, nil
}

func (oidc *GoogleOIDC) VerifyTokenAndGetClaims(ctx context.Context, idToken string) (*Claims, error) {
	verifier := oidc.provider.Verifier(&goOidc.Config{ClientID: oidc.ClientId})

	// Parse and verify ID Token payload.
	parsedToken, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := parsedToken.Claims(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

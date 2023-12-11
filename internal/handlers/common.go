package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tkhq/oauth/internal/helpers/google"
	"github.com/tkhq/oauth/internal/helpers/pages"
	"go.uber.org/zap"
)

/*
HandleHome Function renders the index page when the application index route is called
*/
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pages.IndexPage))
}

func LoginHandler(logger *zap.Logger, oidc *google.GoogleOIDC) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		URL, err := url.Parse(oidc.OauthConfig.Endpoint.AuthURL)
		if err != nil {
			logger.Error("Parse: " + err.Error())
		}
		logger.Info(URL.String())

		parameters := url.Values{}
		parameters.Add("client_id", oidc.OauthConfig.ClientID)
		parameters.Add("scope", strings.Join(oidc.OauthConfig.Scopes, " "))
		parameters.Add("redirect_uri", oidc.OauthConfig.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", oidc.OauthState)
		URL.RawQuery = parameters.Encode()
		url := URL.String()
		logger.Info(url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func CallbackHandler(ctx context.Context, oidc *google.GoogleOIDC, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Google OIDC Callback")

		state := r.FormValue("state")
		if state != oidc.OauthState {
			logger.Info("invalid oauth state, expected " + oidc.OauthState + ", got " + state + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			logger.Warn("Code not found..")
			w.Write([]byte("Code Not Found to provide AccessToken..\n"))
			reason := r.FormValue("error_reason")
			if reason == "user_denied" {
				w.Write([]byte("User has denied Permission.."))
			}
		} else {
			token, err := oidc.GetToken(ctx, code, logger)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			claims, err := oidc.VerifyTokenAndGetClaims(ctx, token)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte("Hello, I'm authenticated\n"))
			w.Write([]byte(fmt.Sprintf("OIDC claims: %+v", claims)))
			return
		}
	}
}

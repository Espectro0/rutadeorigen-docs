package auth

import (
	"fmt"
	"os"

	"github.com/workos/workos-go/v10"
)

type Config struct {
	Client      *workos.Client
	ClientId    string
	RedirectURI string
}

func NewClient() (*Config, error) {
	apiKey := os.Getenv("WORKOS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Missing WORKOS_API_KEY variable.")
	}

	clientId := os.Getenv("WORKOS_CLIENT_ID")
	if clientId == "" {
		return nil, fmt.Errorf("Missing WORKOS_CLIENT_ID variable.")
	}

	redirectURI := os.Getenv("WORKOS_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/callback"
	}

	client := workos.NewClient(apiKey, workos.WithClientID(clientId))

	return &Config{
		Client:      client,
		ClientId:    clientId,
		RedirectURI: redirectURI,
	}, nil
}

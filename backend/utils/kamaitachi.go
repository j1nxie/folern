package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/j1nxie/folern/models"
	"golang.org/x/oauth2"
)

// TODO: figure this out

func EncryptAPIKey(apiKey string) string {
	return apiKey
}

func DecryptAPIKey(encryptedAPIKey string) string {
	return encryptedAPIKey
}

type KamaitachiOAuth2Config struct {
	*oauth2.Config
}

func (c *KamaitachiOAuth2Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	jsonBody := map[string]string{
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"redirect_uri":  c.RedirectURL,
		"grant_type":    "authorization_code",
		"code":          code,
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint.TokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ktResp models.KamaitachiResponse[models.KamaitachiAPITokenResponse]
	if err := json.NewDecoder(resp.Body).Decode(&ktResp); err != nil {
		return nil, err
	}

	if !ktResp.Success {
		return nil, fmt.Errorf("token exchange failed: %s", ktResp.Description)
	}

	return &oauth2.Token{
		AccessToken: ktResp.Body.Token,
		TokenType:   "Bearer",
		Expiry:      time.Time{},
	}, nil
}

func (c *KamaitachiOAuth2Config) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(t),
			Base:   http.DefaultTransport,
		},
	}
}

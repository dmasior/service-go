package turnstile

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	Service struct {
		secret string
		c      *http.Client
	}

	tsResp struct {
		Success bool `json:"success"`
	}
)

const verifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

func NewService(secret string) *Service {
	return &Service{
		secret: secret,
		c:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *Service) CheckToken(ctx context.Context, token string) bool {
	form := url.Values{}
	form.Add("response", token)
	form.Add("secret", t.secret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.c.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var respDecoded tsResp
	if err = json.NewDecoder(resp.Body).Decode(&respDecoded); err != nil {
		return false
	}

	return respDecoded.Success
}

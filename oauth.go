package gokinde

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (cl *Client) OAuth2(ctx context.Context, callbackAuthCode, redirectURL string) (*OAuth2Response, error) {
	values := url.Values{}
	values.Set("client_id", cl.cfg.ClientID)
	values.Set("client_secret", cl.cfg.ClientSecret)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", redirectURL)
	values.Set("code", callbackAuthCode)
	payload := values.Encode()

	postURL := cl.cfg.KindeDomain + OAuth2TokenPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request status: %s, response: %s", resp.Status, string(bodyBytes))
	}

	defer resp.Body.Close()

	respoBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var res OAuth2Response
	err = json.Unmarshal(respoBuf, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w: %s", err, string(respoBuf))
	}

	return &res, nil
}

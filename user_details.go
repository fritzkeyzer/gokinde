package gokinde

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserDetails struct {
	Id             string `json:"id"`
	PreferredEmail string `json:"preferred_email"`
	ProvidedId     string `json:"provided_id"`
	LastName       string `json:"last_name"`
	FirstName      string `json:"first_name"`
	Picture        string `json:"picture"`
}

var DecodeUserDetailsErr = fmt.Errorf("decoding user details")

// GetUserDetails expects a pointer to a var of UserDetails or your own version of it if needed.
func (cl *Client) GetUserDetails(ctx context.Context, jwt string, userPtr any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cl.cfg.KindeDomain+UserProfilePath, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request: %s", resp.Status)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(userPtr)
	if err != nil {
		return fmt.Errorf("%w: %v", DecodeUserDetailsErr, err)
	}

	return nil
}

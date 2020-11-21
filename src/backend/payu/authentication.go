package payugo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // in seconds
	GrantType   string `json:"grant_type"`
}

// Authorize calls authorization endpoint and saves accessToken to perform further calls
func (c *Client) Authorize() error {
	rel, err := url.Parse(authorizeEndpoint)
	if err != nil {
		return err
	}
	u := c.baseURL.ResolveReference(rel)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.cfg.ClientID)
	data.Set("client_secret", c.cfg.ClientSecret)

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return err
	}
	if token := authResponse.AccessToken; token != "" {
		c.accessToken = token
	} else {
		return errors.New("invalid, empty access token")
	}
	return nil
}

func (c Client) authorizeRequest(req *http.Request) error {
	if c.accessToken == "" {
		return ErrNoAccessToken
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	return nil
}

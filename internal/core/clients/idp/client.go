package idp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// IDPClient handles HTTP communication with the Identity Provider
type IDPClient struct {
	httpClient *http.Client
}

// NewIDPClient creates a new IDP client with configured timeout
func NewIDPClient() *IDPClient {
	return &IDPClient{
		httpClient: &http.Client{
			Timeout: constants.IDPRequestTimeout,
		},
	}
}

// ExchangeCodeForToken exchanges an authorization code and
// code_verifier for an access token from the IDP token endpoint.
// This implements the OAuth 2.0 Authorization Code flow with PKCE.
//
// Parameters:
//   - ctx: Context for the HTTP request
//   - code: Authorization code from IDP callback
//   - verifier: PKCE code_verifier used in authorization request
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the IDP token response or an error if exchange fails.
func (c *IDPClient) ExchangeCodeForToken(
	ctx context.Context,
	code string,
	cfg *config.Config,
) (*IDPTokenResponse, error) {
	// Build form-encoded request body
	payload := map[string]string{
		"client_id":     cfg.IDPClientID,
		"client_secret": cfg.IDPClientSecret,
		"code":          code,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Marshal JSON}: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		cfg.IDPTokenURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Create Token Request}: %w", err)
	}

	// Update header to application/json
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Execute Token Request}: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Read Response Body}: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"[IDPClient] {Token Exchange Failed}: status %d, body: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var tokenResp IDPTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("[IDPClient] {Parse Token Response}: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user profile information from the IDP
// userinfo endpoint using an access token.
//
// Parameters:
//   - ctx: Context for the HTTP request
//   - accessToken: Access token from IDP token exchange
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the user information or an error if retrieval fails.
func (c *IDPClient) GetUserInfo(
	ctx context.Context,
	accessToken string,
	cfg *config.Config,
) (*IDPUserInfo, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		cfg.IDPUserinfoURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[IDPClient] {Create UserInfo Request}: %w",
			err,
		)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"[IDPClient] {Execute UserInfo Request}: %w",
			err,
		)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"[IDPClient] {UserInfo Request Failed}: "+
				"status %d, body: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	// Parse response
	var userInfo IDPUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf(
			"[IDPClient] {Parse UserInfo Response}: %w",
			err,
		)
	}

	return &userInfo, nil
}

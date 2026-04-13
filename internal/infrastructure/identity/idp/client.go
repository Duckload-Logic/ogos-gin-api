package idp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// Build request body matching IDP Swagger
	payload := IDPTokenExchangeRequest{
		ClientID:     cfg.IDPClientID,
		ClientSecret: cfg.IDPClientSecret,
		Code:         code,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Marshal JSON}: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/token", cfg.IDPBaseUrl),
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
		fmt.Sprintf("%s/me", cfg.IDPBaseUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[IDPClient] {Create UserInfo Request}: %w",
			err,
		)
	}

	// Set Authorization header
	headerValue := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", headerValue)

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

func (c *IDPClient) RefreshToken(
	ctx context.Context,
	refreshToken string,
	cfg *config.Config,
) (*IDPTokenResponse, error) {
	payload := map[string]string{
		"refresh_token": refreshToken,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Marshal JSON}: %w", err)
	}

	url := fmt.Sprintf("%s/auth/refresh", cfg.IDPBaseUrl)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Create Refresh Request}: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Execute Refresh Request}: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"[IDPClient] {Refresh Failed}: status %d, body: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var tokenResp IDPTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("[IDPClient] {Parse Refresh Response}: %w", err)
	}

	return &tokenResp, nil
}

// ValidateSession checks if the provided session ID is valid by calling
// the IDP's session endpoint with the idp_session cookie.
func (c *IDPClient) ValidateSession(
	ctx context.Context,
	sessionID string,
	cfg *config.Config,
) (*IDPSessionResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/auth/session", cfg.IDPBaseUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Create Session Request}: %w", err)
	}

	// Set the idp_session cookie
	req.AddCookie(&http.Cookie{
		Name:  "idp_session",
		Value: sessionID,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Execute Session Request}: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"[IDPClient] {Session Validation Failed}: status %d, body: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var sessionResp IDPSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("[IDPClient] {Parse Session Response}: %w", err)
	}

	return &sessionResp, nil
}

func (c *IDPClient) Logout(
	ctx context.Context,
	cfg *config.Config,
	accessToken string,
) (*IDPLogoutResponse, error) {
	payload := map[string]string{
		"client_id": cfg.IDPClientID,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Marshal JSON}: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/logout", cfg.IDPBaseUrl),
		bytes.NewReader(jsonBody),
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Create Logout Request}: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[IDPClient] {Execute Logout Request}: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[IDPClient] Error reading logout body: %v", err)
	}

	if resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusNoContent {
		var logoutResp IDPLogoutResponse
		// Attempt to parse JSON if body is not empty
		if len(bodyBytes) > 0 && bodyBytes[0] == '{' {
			_ = json.Unmarshal(bodyBytes, &logoutResp)
		}

		if logoutResp.Message == "" {
			logoutResp.Message = "Logout successful"
		}

		return &logoutResp, nil
	}

	return nil, fmt.Errorf(
		"[IDPClient] {IDP Logout Failed}: status %d, body %s",
		resp.StatusCode,
		string(bodyBytes),
	)
}

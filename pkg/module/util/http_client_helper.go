package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// getOAuth2ClientConfig extracts the OAuth2 client configuration from the provided interface
func GetOAuth2ClientConfig(m interface{}) (*ServiceConfig, error) {
	config := m.(*ServiceConfig)

	if config.TenantURL == "" {
		return nil, fmt.Errorf("tenant_url is not set in the provider configuration")
	}
	if config.ServiceCredClientID == "" {
		return nil, fmt.Errorf("serviceCredClientID is not set in the provider configuration")
	}
	if config.ServiceCredClientSecret == "" {
		return nil, fmt.Errorf("serviceCredClientSecret is not set in the provider configuration")
	}

	return config, nil
}

// getAccessToken fetches an access token using the provided OAuth2 client configuration
func GetAccessToken(config *ServiceConfig) (string, error) {
	tokenURL := fmt.Sprintf("%s%s", config.TenantURL, TokenEndpointPath)
	return FetchAccessToken(tokenURL, config.ServiceCredClientID, config.ServiceCredClientSecret)
}

// fetchAccessToken makes a request to obtain an access token
func FetchAccessToken(tokenURL, clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set(GrantTypeKey, GrantTypeClientCreds)
	data.Set(ClientIDKey, clientID)
	data.Set(ClientSecretKey, clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Printf("[ERROR] Error making token request: %s", err)
		return "", fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Error reading response body: %s", err)
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	accessToken, _, _, _, err := ParseTokenResponse(body)
	if err != nil {
		log.Printf("[ERROR] Error parsing token response: %s", err)
		return "", fmt.Errorf("error parsing token response: %w", err)
	}

	return accessToken, nil
}

// createHTTPRequest creates an HTTP request with the necessary headers
func CreateHTTPRequest(method, url string, body io.Reader, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

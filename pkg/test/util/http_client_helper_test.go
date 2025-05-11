package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	moduleUtil "github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/stretchr/testify/assert"
)

func TestMockTokenHandler(t *testing.T) {
	handler := moduleUtil.MockTokenHandler() // Correctly reference MockTokenHandler
	server := httptest.NewServer(handler)    // Use httptest.NewServer to start a test server
	defer server.Close()

	resp, err := http.Get(server.URL) // Make a test request to the server
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFetchAccessToken(t *testing.T) {
	mockServer := httptest.NewServer(moduleUtil.MockTokenHandler()) // Use MockTokenHandler from mock_token_server.go
	defer mockServer.Close()

	token, err := moduleUtil.FetchAccessToken(mockServer.URL, "mockClientID", "mockClientSecret")

	assert.NoError(t, err)
	assert.Equal(t, "mockAccessToken", token)
}

func TestFetchAccessToken_ErrorScenarios(t *testing.T) {
	// Test invalid token URL
	_, err := moduleUtil.FetchAccessToken(":", "mockClientID", "mockClientSecret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error making token request")

	// Test invalid client credentials
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err = moduleUtil.FetchAccessToken(server.URL, "invalidClientID", "invalidClientSecret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token response")

	// Test invalid response body
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-json"))
	}))
	defer mockServer.Close()

	_, err = moduleUtil.FetchAccessToken(mockServer.URL, "mockClientID", "mockClientSecret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token response")
}

func TestCreateHTTPRequest(t *testing.T) {
	// Test valid request
	req, err := moduleUtil.CreateHTTPRequest("GET", "https://example.com", nil, "mockAccessToken")
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://example.com", req.URL.String())
	assert.Equal(t, "Bearer mockAccessToken", req.Header.Get("Authorization"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

	// Test invalid URL
	_, err = moduleUtil.CreateHTTPRequest("GET", ":", nil, "mockAccessToken")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating request")
}

func TestCreateHTTPRequest_ErrorScenarios(t *testing.T) {
	// Test invalid URL
	_, err := moduleUtil.CreateHTTPRequest("GET", ":", nil, "mockAccessToken")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating request")

	// Test missing access token
	req, err := moduleUtil.CreateHTTPRequest("GET", "https://example.com", nil, "")
	assert.NoError(t, err)
	assert.Empty(t, req.Header.Get("Authorization"))
}

func TestGetOAuth2ClientConfig(t *testing.T) {
	// Test valid configuration
	validConfig := &moduleUtil.ServiceConfig{
		TenantURL:               "https://example.com",
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	}

	config, err := moduleUtil.GetOAuth2ClientConfig(validConfig)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, validConfig, config)
}

func TestGetOAuth2ClientConfig_ErrorScenarios(t *testing.T) {
	// Test missing TenantURL
	_, err := moduleUtil.GetOAuth2ClientConfig(&moduleUtil.ServiceConfig{
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant_url is not set in the provider configuration")

	// Test missing ServiceCredClientID
	_, err = moduleUtil.GetOAuth2ClientConfig(&moduleUtil.ServiceConfig{
		TenantURL:               "https://example.com",
		ServiceCredClientSecret: "mockClientSecret",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "serviceCredClientID is not set in the provider configuration")

	// Test missing ServiceCredClientSecret
	_, err = moduleUtil.GetOAuth2ClientConfig(&moduleUtil.ServiceConfig{
		TenantURL:           "https://example.com",
		ServiceCredClientID: "mockClientID",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "serviceCredClientSecret is not set in the provider configuration")
}

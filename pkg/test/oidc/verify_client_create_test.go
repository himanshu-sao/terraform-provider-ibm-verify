package oidc_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/oidc"
	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestResourceIBMVerifyOAuth2ClientCreate(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock Server: Received request - Method: %s, URL: %s", r.Method, r.URL.String())
		if r.URL.Path == "/v1.0/endpoint/default/token" { // Ensure this matches the expected token endpoint
			util.MockTokenHandler()(w, r)
		} else if r.URL.Path == "/oauth2/register" {
			// Mock response for client creation
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			response := `{"client_id": "mockClientID"}`
			log.Printf("Mock Server: Sending response - %s", response)
			w.Write([]byte(response))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Mock provider configuration with all required fields
	mockConfig := &util.ServiceConfig{
		TenantURL:               mockServer.URL,
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	}

	resource := oidc.ResourceIBMVerifyOAuth2Client()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"client_name":   "TestClient",
		"redirect_uris": []interface{}{"https://example.com/callback"}, // Add redirect_uris
	})
	err := resource.Create(d, mockConfig)

	// Log the resource data for debugging
	log.Printf("Resource Data: ID: %s, Error: %v", d.Id(), err)

	assert.NoError(t, err)
	assert.Equal(t, "mockClientID", d.Id())
}

func TestResourceIBMVerifyOAuth2ClientCreate_InvalidTokenResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock Server: Received request - Method: %s, URL: %s", r.Method, r.URL.String())
		if r.URL.Path == "/v1.0/endpoint/default/token" {
			// Return an invalid JSON response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-json`))
		} else if r.URL.Path == "/oauth2/register" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"client_id": "mockClientID"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	mockConfig := &util.ServiceConfig{
		TenantURL:               mockServer.URL,
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	}

	resource := oidc.ResourceIBMVerifyOAuth2Client()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"client_name":   "TestClient",
		"redirect_uris": []interface{}{"https://example.com/callback"},
	})
	err := resource.Create(d, mockConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing JSON response")
}

func TestResourceIBMVerifyOAuth2ClientCreate_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	mockConfig := &util.ServiceConfig{
		TenantURL:               mockServer.URL,
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	}

	resource := oidc.ResourceIBMVerifyOAuth2Client()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"client_name":   "TestClient",
		"redirect_uris": []interface{}{"https://example.com/callback"},
	})

	err := resource.Create(d, mockConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token response") // Updated assertion
}

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

func TestResourceIBMVerifyOAuth2ClientRead(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock Server: Received request - Method: %s, URL: %s", r.Method, r.URL.String())
		if r.URL.Path == "/v1.0/endpoint/default/token" { // Ensure this matches the expected token endpoint
			util.MockTokenHandler()(w, r)
		} else if r.URL.Path == "/oauth2/register/mockClientID" {
			// Mock response for client read
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"client_id": "mockClientID", "client_name": "TestClient"}`
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
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("mockClientID") // Set the ID for the resource

	err := resource.Read(d, mockConfig)

	// Log the resource data for debugging
	log.Printf("Resource Data: ID: %s, Client Name: %s, Error: %v", d.Id(), d.Get("client_name"), err)

	assert.NoError(t, err)
	assert.Equal(t, "mockClientID", d.Id())
	assert.Equal(t, "TestClient", d.Get("client_name"))
}

func TestResourceIBMVerifyOAuth2ClientRead_UnexpectedStatusCode(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock Server: Received request - Method: %s, URL: %s", r.Method, r.URL.String())
		if r.URL.Path == "/v1.0/endpoint/default/token" {
			util.MockTokenHandler()(w, r)
		} else if r.URL.Path == "/oauth2/register/mockClientID" {
			// Return an unexpected status code
			w.WriteHeader(http.StatusInternalServerError)
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
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("mockClientID")

	err := resource.Read(d, mockConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestResourceIBMVerifyOAuth2ClientRead_InvalidJSONResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid-json`))
	}))
	defer mockServer.Close()

	mockConfig := &util.ServiceConfig{
		TenantURL:               mockServer.URL,
		ServiceCredClientID:     "mockClientID",
		ServiceCredClientSecret: "mockClientSecret",
	}

	resource := oidc.ResourceIBMVerifyOAuth2Client()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("mockClientID")

	err := resource.Read(d, mockConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token response") // Updated assertion
}

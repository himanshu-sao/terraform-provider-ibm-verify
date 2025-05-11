package oidc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/oidc"
	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestResourceIBMVerifyOAuth2ClientUpdate(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1.0/endpoint/default/token" {
			util.MockTokenHandler()(w, r)
		} else if r.URL.Path == "/oauth2/register/mockClientID" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"client_id": "mockClientID", "client_name": "UpdatedClient"}`))
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
		"client_name":   "UpdatedClient",
		"redirect_uris": []interface{}{"https://example.com/callback"}, // Ensure redirect_uris is included
	})
	d.SetId("mockClientID")
	err := resource.Update(d, mockConfig)

	assert.NoError(t, err)
	assert.Equal(t, "mockClientID", d.Id())
	assert.Equal(t, "UpdatedClient", d.Get("client_name"))
}

func TestResourceIBMVerifyOAuth2ClientUpdate_InvalidPayload(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a timeout error by closing the connection without responding
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
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
		"client_name": 123, // Invalid type for client_name
	})
	d.SetId("mockClientID")

	err := resource.Update(d, mockConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error making token request")
}

func TestResourceIBMVerifyOAuth2ClientDelete(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1.0/endpoint/default/token" {
			util.MockTokenHandler()(w, r)
		} else if r.URL.Path == "/oauth2/register/mockClientID" {
			w.WriteHeader(http.StatusNoContent)
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
	err := resource.Delete(d, mockConfig)

	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

func TestResourceIBMVerifyOAuth2ClientDelete_HTTPError(t *testing.T) {
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
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("mockClientID")

	err := resource.Delete(d, mockConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing token response")
}

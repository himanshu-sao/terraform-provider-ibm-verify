package util

import (
	"encoding/json"
	"testing"

	moduleUtil "github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/stretchr/testify/assert"
)

func TestParseTokenResponse(t *testing.T) {
	mockResponse := map[string]interface{}{
		"access_token": "mockAccessToken",
		"grant_id":     "mockGrantID",
		"token_type":   "Bearer",
		"expires_in":   3600,
	}
	body, _ := json.Marshal(mockResponse)

	accessToken, grantID, tokenType, expiresIn, err := moduleUtil.ParseTokenResponse(body)

	assert.NoError(t, err)
	assert.Equal(t, "mockAccessToken", accessToken)
	assert.Equal(t, "mockGrantID", grantID)
	assert.Equal(t, "Bearer", tokenType)
	assert.Equal(t, 3600, expiresIn)
}

func TestParseTokenResponse_ErrorScenarios(t *testing.T) {
	// Test invalid JSON
	_, _, _, _, err := moduleUtil.ParseTokenResponse([]byte("invalid-json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing JSON response")

	// Test missing fields
	mockResponse := map[string]interface{}{
		"access_token": "mockAccessToken",
	}
	body, _ := json.Marshal(mockResponse)

	_, _, _, _, err = moduleUtil.ParseTokenResponse(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grant_id not found in response")
}

func TestParseTokenResponse_EmptyResponse(t *testing.T) {
	// Test empty response
	_, _, _, _, err := moduleUtil.ParseTokenResponse([]byte("{}"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access_token not found in response")
}

func TestParseTokenResponse_InvalidDataType(t *testing.T) {
	// Test invalid data type for expires_in
	mockResponse := map[string]interface{}{
		"access_token": "mockAccessToken",
		"grant_id":     "mockGrantID",
		"token_type":   "Bearer",
		"expires_in":   "invalid",
	}
	body, _ := json.Marshal(mockResponse)

	_, _, _, _, err := moduleUtil.ParseTokenResponse(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expires_in not found in response")
}

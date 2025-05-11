package oidc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// handleHTTPResponse checks if the HTTP response status code matches the expected status
func handleHTTPResponse(resp *http.Response, expectedStatus int) error {
	if resp.StatusCode != expectedStatus {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] Unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// buildPayloadFromSchema constructs a payload from the Terraform schema data
func buildPayloadFromSchema(resourceData *schema.ResourceData) (map[string]interface{}, error) {
	payload := make(map[string]interface{})

	if jsonPath, ok := resourceData.GetOk(util.JsonConfigPathKey); ok {
		log.Printf("[DEBUG] Using JSON config file: %s", jsonPath.(string))
		jsonFile, err := os.ReadFile(jsonPath.(string))
		if err != nil {
			log.Printf("[ERROR] Failed to read JSON file: %s", err)
			return nil, fmt.Errorf("error reading JSON file: %w", err)
		}

		err = json.Unmarshal(jsonFile, &payload)
		if err != nil {
			log.Printf("[ERROR] Failed to parse JSON file: %s", err)
			return nil, fmt.Errorf("error parsing JSON file: %w", err)
		}
	} else {
		for _, key := range util.SchemaFields {
			if !util.ExcludedFieldsFromPayload[key] {
				if v, ok := resourceData.GetOk(key); ok {
					payload[key] = v
				}
			}
		}
	}

	// Log the payload before removing excluded fields
	log.Printf("[ERROR] Payload before removing excluded fields: %+v", payload)

	// Remove excluded fields from payload
	for key := range util.ExcludedFieldsFromPayload {
		delete(payload, key)
	}

	// Log the payload after removing excluded fields
	log.Printf("[ERROR] Payload after removing excluded fields: %+v", payload)

	if err := validateJSONPayload(payload); err != nil {
		log.Printf("[ERROR] Invalid payload: %s", err)
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	return payload, nil
}

// validateJSONPayload checks if the payload contains all required fields and if they are of the correct type
func validateJSONPayload(payload map[string]interface{}) error {
	required := []string{"redirect_uris", "client_name"}
	for _, field := range required {
		if _, ok := payload[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	if redirectURIs, ok := payload["redirect_uris"].([]interface{}); ok {
		if len(redirectURIs) == 0 {
			return fmt.Errorf("redirect_uris must not be empty")
		}
		for _, uri := range redirectURIs {
			if _, ok := uri.(string); !ok {
				return fmt.Errorf("redirect_uris must be an array of strings")
			}
		}
	} else {
		return fmt.Errorf("redirect_uris must be an array")
	}

	if _, ok := payload["client_name"].(string); !ok {
		return fmt.Errorf("client_name must be a string")
	}

	return nil
}

// sendRequest sends an HTTP request with the given method, URL, payload, and access token
func sendRequest(method, url string, payload map[string]interface{}, accessToken string) (*http.Response, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request body: %s", err)
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	client := &http.Client{}
	req, err := util.CreateHTTPRequest(method, url, bytes.NewReader(requestBody), accessToken)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Sending request to %s", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed to make request: %s", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}

// parseAndSetOAuth2ClientData parses the response body and sets the data in the Terraform schema
func parseAndSetOAuth2ClientData(resourceData *schema.ResourceData, body io.Reader) error {
	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding JSON response: %w", err)
	}

	for k, v := range result {
		if err := resourceData.Set(k, v); err != nil {
			return fmt.Errorf("error setting %s: %w", k, err)
		}
	}

	return nil
}

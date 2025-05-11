package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProfileFlowRead(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Get the flow ID from the resource data
	flowID := d.Id()
	if flowID == "" {
		return fmt.Errorf("flow ID is not set")
	}

	// Construct the GET URL
	url := fmt.Sprintf("%s/profile/config/v3.0/flows/%s", config.TenantURL, flowID)

	// Create the GET request
	req, err := util.CreateHTTPRequest("GET", url, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusNotFound {
		// If the flow is not found, remove it from the state
		d.SetId("")
		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var flow map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&flow); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Set the resource data
	if err := d.Set("name", flow["name"]); err != nil {
		return fmt.Errorf("failed to set name: %w", err)
	}
	if err := d.Set("type", flow["type"]); err != nil {
		return fmt.Errorf("failed to set type: %w", err)
	}
	if err := d.Set("enabled", flow["enabled"]); err != nil {
		return fmt.Errorf("failed to set enabled: %w", err)
	}

	return nil
}

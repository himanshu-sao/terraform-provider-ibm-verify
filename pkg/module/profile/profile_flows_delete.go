package profile

import (
	"fmt"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProfileFlowDelete(d *schema.ResourceData, m interface{}) error {
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

	// Construct the DELETE URL
	url := fmt.Sprintf("%s/profile/config/v3.0/flows/%s", config.TenantURL, flowID)

	// Create the DELETE request
	req, err := util.CreateHTTPRequest("DELETE", url, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Successfully deleted the flow
	d.SetId("") // Clear the resource ID
	return nil
}

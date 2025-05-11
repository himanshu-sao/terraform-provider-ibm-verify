package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceProfileFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileFlowCreate,
		Read:   dataSourceProfileFlowsRead,
		Delete: resourceProfileFlowDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceProfileFlowCreate(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	payload := map[string]interface{}{
		"name":    d.Get("name").(string),
		"type":    d.Get("type").(string),
		"enabled": d.Get("enabled").(bool),
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/profile/config/v3.0/flows", config.TenantURL)
	req, err := util.CreateHTTPRequest("POST", url, bytes.NewReader(requestBody), accessToken)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	flowID, ok := responseBody["id"].(string)
	if !ok {
		return fmt.Errorf("failed to get flow ID from response")
	}

	d.SetId(flowID)
	return nil
}

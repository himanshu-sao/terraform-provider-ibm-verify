package profile

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceProfileFlows() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProfileFlowsRead,

		Schema: map[string]*schema.Schema{
			"flows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"published": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProfileFlowsRead(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/profile/config/v3.0/flows", config.TenantURL)
	req, err := util.CreateHTTPRequest("GET", url, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Flows []map[string]interface{} `json:"flows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := d.Set("flows", result.Flows); err != nil {
		return fmt.Errorf("failed to set flows: %w", err)
	}

	d.SetId("profile_flows")
	log.Printf("[DEBUG] Successfully fetched profile flows")
	return nil
}

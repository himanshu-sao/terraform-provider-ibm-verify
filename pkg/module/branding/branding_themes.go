package branding

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBrandingThemes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBrandingThemesRead,

		Schema: map[string]*schema.Schema{
			"themes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBrandingThemesRead(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s%s", config.TenantURL, BrandingThemesURL)
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
		ThemeRegistrations []map[string]interface{} `json:"themeRegistrations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	themes := make([]map[string]interface{}, len(result.ThemeRegistrations))
	for i, theme := range result.ThemeRegistrations {
		themes[i] = map[string]interface{}{
			"name":        theme["name"],
			"description": theme["description"],
			"id":          theme["id"],
		}
	}

	if err := d.Set("themes", themes); err != nil {
		return fmt.Errorf("failed to set themes: %w", err)
	}

	d.SetId("branding_themes")
	log.Printf("[DEBUG] Successfully fetched branding themes")
	return nil
}

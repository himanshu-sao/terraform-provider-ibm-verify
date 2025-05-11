package branding

import (
	"fmt"
	"io"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMVerifyBrandingElement() *schema.Resource {
	return &schema.Resource{
		Create: resourceBrandingThemeCreate,
		Read:   resourceBrandingThemeRead,
		Update: resourceBrandingThemeUpdate,
		Delete: resourceBrandingThemeDelete,

		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The absolute path to the theme file (e.g., Template.zip).",
			},
			"theme_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the branding theme.",
			},
			"theme_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the branding theme.",
			},
		},
	}
}

func resourceBrandingThemeRead(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	themeID := d.Id()
	url := fmt.Sprintf("%s%s/%s", config.TenantURL, BrandingThemesURL, themeID)
	req, err := util.CreateHTTPRequest("GET", url, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response (if needed)
	// Example: d.Set("theme_name", "parsed_name")

	return nil
}

func resourceBrandingThemeDelete(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	themeID := d.Id()
	url := fmt.Sprintf("%s%s/%s", config.TenantURL, BrandingThemesURL, themeID)
	req, err := util.CreateHTTPRequest("DELETE", url, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	d.SetId("")
	return nil
}

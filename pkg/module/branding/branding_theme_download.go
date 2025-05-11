package branding

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBrandingThemeDownload() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBrandingThemeDownloadRead,

		Schema: map[string]*schema.Schema{
			"theme_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "The ID of the branding theme to download. Defaults to 'default'.",
			},
			"download_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The absolute path where the theme file will be downloaded.",
			},
		},
	}
}

func dataSourceBrandingThemeDownloadRead(d *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	themeID := d.Get("theme_id").(string)
	downloadPath := d.Get("download_path").(string)

	url := fmt.Sprintf("%s%s/%s", config.TenantURL, BrandingThemesURL, themeID)
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

	// Ensure the response is of type application/octet-stream
	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		return fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	// Create the file at the specified download path
	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %w", downloadPath, err)
	}
	defer file.Close()

	// Write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	d.SetId(themeID)
	return nil
}

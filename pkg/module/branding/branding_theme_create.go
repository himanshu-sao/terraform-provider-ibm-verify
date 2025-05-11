package branding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Handles the creation of a branding theme
func resourceBrandingThemeCreate(d *schema.ResourceData, m interface{}) error {
	return createOrUpdateTheme(d, m, "POST", "")
}

// Handles the update of a branding theme
func resourceBrandingThemeUpdate(d *schema.ResourceData, m interface{}) error {
	themeID := d.Id()
	return createOrUpdateTheme(d, m, "PUT", themeID)
}

// Shared function for creating or updating a branding theme
func createOrUpdateTheme(d *schema.ResourceData, m interface{}, method, themeID string) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 client config: %w", err)
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	filePath := d.Get("file_path").(string)
	themeName := d.Get("theme_name").(string)
	themeDescription := d.Get("theme_description").(string)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Create a multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the form
	fileWriter, err := writer.CreateFormFile("files", filePath)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add the configuration JSON to the form
	configuration := map[string]string{
		"name":        themeName,
		"description": themeDescription,
	}
	configurationJSON, err := json.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}
	if err := writer.WriteField("configuration", string(configurationJSON)); err != nil {
		return fmt.Errorf("failed to write configuration field: %w", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s%s", config.TenantURL, BrandingThemesURL)
	if method == "PUT" {
		url = fmt.Sprintf("%s%s/%s", config.TenantURL, BrandingThemesURL, themeID)
	}
	req, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response to get the theme ID
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if method == "POST" {
		themeID, ok := responseBody["id"].(string)
		if !ok {
			return fmt.Errorf("failed to get theme ID from response")
		}
		d.SetId(themeID)
	}

	return nil
}

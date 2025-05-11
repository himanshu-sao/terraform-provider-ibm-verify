package branding

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Theme struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type ListThemesResponse struct {
	Count              int     `json:"count"`
	Limit              int     `json:"limit"`
	Page               int     `json:"page"`
	Total              int     `json:"total"`
	ThemeRegistrations []Theme `json:"themeRegistrations"`
}

// resourceIBMVerifyBrandingElementListAllThemes fetches all themes from the API.
func resourceIBMVerifyBrandingElementListAllThemes(authToken, tenantHost string) (*ListThemesResponse, error) {
	// Construct the API endpoint dynamically using the constant
	url := fmt.Sprintf("https://%s%s", tenantHost, BrandingThemesURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the Authorization header
	req.Header.Add("Authorization", "Bearer "+authToken)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var listThemesResponse ListThemesResponse
	if err := json.Unmarshal(body, &listThemesResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &listThemesResponse, nil
}

// resourceIBMVerifyBrandingElementGetDefaultTheme fetches the default theme and saves it to the specified file.
func resourceIBMVerifyBrandingElementGetDefaultTheme(authToken, tenantHost, outputFilePath string) error {
	// Construct the API endpoint dynamically using the constant
	url := fmt.Sprintf("https://%s%s", tenantHost, DefaultBrandingThemeURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add the Authorization header
	req.Header.Add("Authorization", "Bearer "+authToken)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write the response body to the file
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	return nil
}

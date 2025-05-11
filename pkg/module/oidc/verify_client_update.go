package oidc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMVerifyOAuth2ClientUpdate(resourceData *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return err
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return err
	}

	clientID := resourceData.Id()
	log.Printf("[DEBUG] Updating OAuth2 client with ID: %s. Tenant URL: %s", clientID, config.TenantURL)

	payload, err := buildPayloadFromSchema(resourceData)
	if err != nil {
		return err
	}

	resp, err := sendRequest("PUT", fmt.Sprintf("%s%s/%s", config.TenantURL, util.OauthClientEndpoint, clientID), payload, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := handleHTTPResponse(resp, http.StatusOK); err != nil {
		return err
	}

	if err := parseAndSetOAuth2ClientData(resourceData, resp.Body); err != nil {
		log.Printf("[ERROR] Failed to parse response: %s", err)
		return fmt.Errorf("error parsing response: %w", err)
	}

	log.Printf("[DEBUG] Successfully updated OAuth2 client with ID: %s", resourceData.Id())
	return nil
}

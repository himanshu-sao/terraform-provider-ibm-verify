package oidc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMVerifyOAuth2ClientCreate(resourceData *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return err
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating OAuth2 client. Tenant URL: %s", config.TenantURL)

	payload, err := buildPayloadFromSchema(resourceData)
	if err != nil {
		return err
	}

	resp, err := sendRequest("POST", fmt.Sprintf("%s%s", config.TenantURL, util.OauthClientEndpoint), payload, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := handleHTTPResponse(resp, http.StatusCreated); err != nil {
		return err
	}

	if err := parseAndSetOAuth2ClientData(resourceData, resp.Body); err != nil {
		log.Printf("[ERROR] Failed to parse response: %s", err)
		return fmt.Errorf("error parsing response: %w", err)
	}

	resourceData.SetId(resourceData.Get("client_id").(string))
	log.Printf("[DEBUG] Successfully created OAuth2 client with ID: %s", resourceData.Id())
	return nil
}

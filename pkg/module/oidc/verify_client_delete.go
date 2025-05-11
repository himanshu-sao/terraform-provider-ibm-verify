package oidc

import (
	"fmt"
	"net/http"

	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMVerifyOAuth2ClientDelete(resourceData *schema.ResourceData, m interface{}) error {
	config, err := util.GetOAuth2ClientConfig(m)
	if err != nil {
		return err
	}

	accessToken, err := util.GetAccessToken(config)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := util.CreateHTTPRequest("DELETE", fmt.Sprintf("%s%s/%s", config.TenantURL, util.OauthClientEndpoint, resourceData.Id()), nil, accessToken)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if err := handleHTTPResponse(resp, http.StatusNoContent); err != nil {
		// Check for OK status as well, in case the API returns OK for successful deletion
		if resp.StatusCode != http.StatusOK {
			return err
		}
	}

	resourceData.SetId("")
	return nil
}

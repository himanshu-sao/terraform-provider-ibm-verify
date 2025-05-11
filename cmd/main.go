package main

import (
	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/branding"
	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/oidc"
	"github.com/IBM-Verify/terraform-provider-ibm-verify/pkg/module/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IBM_VERIFY_TENANT_URL", nil),
				Description: "The IBM Verify tenant URL",
			},
			"service_cred_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IBM_VERIFY_SERVICE_CRED_CLIENT_ID", nil),
				Description: "The service credential client ID",
			},
			"service_cred_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IBM_VERIFY_SERVICE_CRED_CLIENT_SECRET", nil),
				Description: "The service credential client secret",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ibmverify_oauth2_client":    oidc.ResourceIBMVerifyOAuth2Client(),
			"ibmverify_branding_element": branding.ResourceIBMVerifyBrandingElement(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"ibmverify_oauth2_client": oidc.DataSourceIBMVerifyOAuth2Client(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(resourceData *schema.ResourceData) (interface{}, error) {
	config := util.ServiceConfig{
		TenantURL:               resourceData.Get("tenant_url").(string),
		ServiceCredClientID:     resourceData.Get("service_cred_client_id").(string),
		ServiceCredClientSecret: resourceData.Get("service_cred_client_secret").(string),
	}
	return &config, nil
}

terraform {
  #required_providers {
  #  ibmverify = {
  #    source  = "registry.terraform.io/local/ibmverify"
  #    version = "0.0.1"
  #  }
  #  time = {
  #    source = "hashicorp/time"
  #  }
  #}
}

# Variables
variable "tenant_url" {
  type        = string
  description = "The tenant URL for IBM Verify"
}

variable "service_cred_client_id" {
  type        = string
  description = "The service credential client ID"
}

variable "service_cred_client_secret" {
  type        = string
  description = "The service credential client secret"
  sensitive   = true
}


provider "ibmverify" {
  tenant_url                 = var.tenant_url
  service_cred_client_id     = var.service_cred_client_id
  service_cred_client_secret = var.service_cred_client_secret
}

# Main resource configuration
resource "ibmverify_oauth2_client" "verify_oidc_client" {
  client_name                            = "Sorceress Supreme"
  access_policy                          = "default"
  all_users_entitled                     = true
  authorization_encrypted_response_alg   = "none"
  authorization_encrypted_response_enc   = "none"
  authorization_signed_response_alg      = "RS256"
  client_secret_expires_at               = 0
  consent_action                         = "always_prompt"
  enforce_pkce                           = true
  grant_types                            = ["authorization_code", "refresh_token"]
  id_token_signed_response_alg           = "RS256"
  redirect_uris                          = ["https://cleanly-renewed-iguana.ngrok-free.app/callback"]
  request_object_check_expiry            = true
  request_object_lifetime                = 1800
  request_object_parameters_only         = false
  request_object_signing_alg             = "RS256"
  require_pushed_authorization_requests  = false
  response_types                         = ["none", "code"]
  restrict_api_entitlements              = false
  token_endpoint_auth_method             = "default"
  token_endpoint_auth_single_use_jti     = false
}

# Output blocks
output "oauth2_client_id" {
  description = "The client ID of the created OAuth2 client"
  value       = ibmverify_oauth2_client.verify_oidc_client.client_id
}

output "oauth2_client_secret" {
  description = "The client secret of the created OAuth2 client"
  value       = ibmverify_oauth2_client.verify_oidc_client.client_secret
  sensitive   = true
}

output "oauth2_registration_access_token" {
  description = "The registration access token of the created OAuth2 client"
  value       = ibmverify_oauth2_client.verify_oidc_client.registration_access_token
  sensitive   = true
}

output "oauth2_registration_client_uri" {
  description = "The registration client URI of the created OAuth2 client"
  value       = ibmverify_oauth2_client.verify_oidc_client.registration_client_uri
}
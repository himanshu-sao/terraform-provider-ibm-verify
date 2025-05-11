terraform {
  required_providers {
    ibmverify = {
      source  = "registry.terraform.io/local/ibmverify"
      version = "1.0.0"
    }
    time = {
      source = "hashicorp/time"
    }
  }
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

# Data block to fetch branding themes
data "ibmverify_branding_themes" "themes" {}

data "ibmverify_branding_theme_download" "theme_download" {
  theme_id      = "default" # Optional, defaults to "default"
  download_path = "/Users/Shared/themes/default_theme.zip" # Replace with your desired path
}

# Output block to display branding themes
output "branding_themes" {
  description = "List of branding themes"
  value       = data.ibmverify_branding_themes.themes
}

output "downloaded_theme_path" {
  description = "The path where the branding theme was downloaded"
  value       = data.ibmverify_branding_theme_download.theme_download.download_path
}

resource "ibmverify_branding_theme" "test_theme" {
  file_path         = "/Users/Shared/themes/default_theme.zip" # Replace with your file path
  theme_name        = "TestTheme"
  theme_description = "This is a test theme for delete functionality"
}

output "test_theme_id" {
  description = "The ID of the test branding theme"
  value       = ibmverify_branding_theme.test_theme.id
}
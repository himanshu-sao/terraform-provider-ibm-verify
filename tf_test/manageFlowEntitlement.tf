terraform {
  required_providers {
    ibmverify = {
      source  = "registry.terraform.io/local/ibmverify"
      version = "0.0.1"
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

resource "ibmverify_update_entitlement" "test" {
  entitlement   = "manageFlows"
  tenant_url    = var.tenant_url
  client_id     = var.service_cred_client_id
  client_secret = var.service_cred_client_secret
}

output "update_entitlement_status" {
  value = "Entitlement update executed for client ID: ${ibmverify_update_entitlement.test.id}"
}
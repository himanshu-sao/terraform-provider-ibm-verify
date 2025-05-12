terraform {
  required_providers {
    ibmverify = {
      source  = "registry.terraform.io/himanshu-sao/ibmverify"
      #source = "registry.terraform.io/local/ibmverify"
      version = ">= 0.0.1"
    }
  }
  required_version = ">= 0.0.1"
}
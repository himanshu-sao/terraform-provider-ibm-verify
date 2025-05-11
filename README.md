# IBM Verify Terraform Provider

The **IBM Verify Terraform Provider** enables seamless integration with **IBM Security Verify** to manage **OIDC applications** and **branding themes** programmatically using Terraform. This provider simplifies the automation of IBM Verify configurations, allowing you to manage your identity and access management (IAM) infrastructure as code.

---

## Table of Contents

1. [Overview](#overview)
2. [Requirements](#requirements)
3. [Installation](#installation)
   - [Using the Makefile](#using-the-makefile)
   - [Building the Provider](#building-the-provider)
   - [Docker Image for the Provider](#docker-image-for-the-provider)
   - [Download the Provider from the Terraform Registry](#download-the-provider-from-the-terraform-registry)
   - [Download the Provider Manually](#download-the-provider-manually)
4. [Provider Configuration](#provider-configuration)
5. [Resources](#resources)
6. [Data Sources](#data-sources)
7. [Testing](#testing)
8. [Contributing](#contributing)
9. [References](#references)
10. [License](#license)

---

## Overview

The IBM Verify Terraform Provider allows you to manage **OIDC applications** and **branding themes** in IBM Security Verify. It provides resources and data sources to automate the configuration of your identity and access management (IAM) infrastructure.

---

## Requirements

- **Terraform**: `>= 0.0.1`
- **Go**: `>= 1.20`
- **IBM Verify Tenant**: A valid IBM Verify tenant URL.
- **Service Credentials**: Client ID and Client Secret for authentication.

---

## Installation

### Using the Makefile

A `Makefile` is included in the repository to simplify common tasks such as building, testing, and packaging the provider. Below are the available commands:

#### Commands:
- **Build the plugin**:
  ```bash
  make build
  ```
  This command runs the `buildPlugin.sh` script to compile the provider.

- **Run tests**:
  ```bash
  make test
  ```
  This command runs the `runTests.sh` script to execute unit tests.

- **Package the plugin**:
  ```bash
  make package
  ```
  This command packages the provider for distribution.

- **Install the plugin locally**:
  ```bash
  make install
  ```
  This command installs the provider locally in the Terraform plugins directory.

- **Validate the plugin with Terraform**:
  ```bash
  make validate
  ```
  This command runs `terraform init` and `terraform validate` to ensure the provider works correctly.

- **Run the full pipeline**:
  ```bash
  make pipeline
  ```
  This command runs the full pipeline: clean, build, test, package, and validate.

---

### Building the Provider

If you prefer not to use the `Makefile`, you can manually build the provider using the following steps:

Clone the repository:
```bash
git clone https://github.com/IBM-Verify/terraform-provider-ibm-verify.git
cd terraform-provider-ibm-verify
```

Build the provider:
```bash
./bin/buildPlugin.sh
```

---

### Download the Provider from the Terraform Registry

1. [Download and install Terraform for your system](https://www.terraform.io/intro/getting-started/install.html).

2. Create a `versions.tf` file in your Terraform module folder and add the following block:
```hcl
terraform {
  required_providers {
    ibmverify = {
      source  = "IBM-Verify/ibmverify"
      version = "<provider version>"
    }
  }
}
```

3. Run `terraform init` to fetch the IBM Verify provider plugin for Terraform from the Terraform Registry.

---

### Download the Provider Manually

1. [Download the IBM Verify provider plugin for Terraform](https://github.com/IBM-Verify/terraform-provider-ibm-verify).

2. Unzip the release archive to extract the plugin binary.

3. Move the binary into the Terraform plugins directory:
   - Linux/Unix/OS X: `~/.terraform.d/plugins`
   - Windows: `%APPDATA%\terraform.d\plugins`

4. Add the provider to your Terraform configuration file:
```hcl
provider "ibmverify" {}
```

---

## Provider Configuration

To use the provider, add the following configuration to your Terraform file:

```hcl
terraform {
  required_providers {
    ibmverify = {
      source  = "IBM-Verify/ibmverify"
      version = "0.0.1"
    }
  }
}

provider "ibmverify" {
  tenant_url                 = var.tenant_url
  service_cred_client_id     = var.service_cred_client_id
  service_cred_client_secret = var.service_cred_client_secret
}
```

### Variables
- `tenant_url`: The tenant URL for IBM Verify.
- `service_cred_client_id`: The client ID for service credentials.
- `service_cred_client_secret`: The client secret for service credentials.

---

## Resources

### OIDC Application Management

#### Resource: `ibmverify_oauth2_client`

This resource allows you to create and manage OAuth2/OIDC clients in IBM Verify.

#### Example Configuration

```hcl
resource "ibmverify_oauth2_client" "verify_oidc_client" {
  client_name                            = "Example OIDC Client"
  access_policy                          = "default"
  all_users_entitled                     = true
  authorization_signed_response_alg      = "RS256"
  client_secret_expires_at               = 0
  consent_action                         = "always_prompt"
  enforce_pkce                           = true
  grant_types                            = ["authorization_code", "refresh_token"]
  redirect_uris                          = ["https://example.com/callback"]
  response_types                         = ["code"]
  token_endpoint_auth_method             = "client_secret_basic"
}
```

---

## Data Sources

### Fetch All Branding Themes
```hcl
data "ibmverify_branding_themes" "themes" {}

output "branding_themes" {
  description = "List of branding themes"
  value       = data.ibmverify_branding_themes.themes
}
```

### Download a Branding Theme
```hcl
data "ibmverify_branding_theme_download" "theme_download" {
  theme_id      = "default"
  download_path = "/path/to/downloaded_theme.zip"
}
```

---

## Unit Testing

Run the following script to execute tests for the provider:

```bash
make test
```

This script runs tests for the following modules:
- `oidc`
- `branding`
- `util`

---

## Contributing

We welcome contributions to this project! Follow these steps to contribute:

1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add your message here"
   ```
4. Push to the branch:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Open a pull request.

---

## References

- [IBM Security Verify Documentation](https://www.ibm.com/docs/en/security-verify)
- [Managing Branding Themes](https://www.ibm.com/docs/en/security-verify?topic=branding-managing-themes)
- [Terraform Documentation](https://www.terraform.io/docs)

---

## License

This project is licensed under the [Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/).
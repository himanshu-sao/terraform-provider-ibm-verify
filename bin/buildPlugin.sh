#!/bin/bash

# Variables for paths
PLUGIN_NAME="terraform-provider-ibmverify"
PLUGIN_VERSION="1.0.0"
PLUGIN_OS_ARCH="darwin_amd64"
PLUGIN_LOCAL_PATH="compiled_package/$PLUGIN_NAME"
PLUGIN_INSTALL_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/local/ibmverify/$PLUGIN_VERSION/$PLUGIN_OS_ARCH"
PLUGIN_INSTALL_PATH="$PLUGIN_INSTALL_DIR/${PLUGIN_NAME}_v${PLUGIN_VERSION}"

echo "Cleaning up the old build"
rm -rf go.mod go.sum
rm -rf .terraform*

mkdir -p compiled_package

# Backup existing plugin if it exists
if [ -e "$PLUGIN_INSTALL_PATH" ]; then
    mv -v "$PLUGIN_INSTALL_PATH" "${PLUGIN_INSTALL_PATH}_bkp"
else
    echo "No existing plugin found at $PLUGIN_INSTALL_PATH. Skipping backup."
fi
sleep 0.5

# Remove old compiled package
rm -rv "$PLUGIN_LOCAL_PATH"

echo "$(tput setaf 3) #####   Building the plugin   #####  $(tput sgr0)"
go mod init github.com/IBM-Verify/terraform-provider-ibm-verify && go mod tidy && go build -o "$PLUGIN_LOCAL_PATH" ./cmd

# Ensure the target directory exists
mkdir -p "$PLUGIN_INSTALL_DIR"

# Copy the compiled plugin to the target directory
cp -v "$PLUGIN_LOCAL_PATH" "$PLUGIN_INSTALL_PATH"

# Set executable permissions for the plugin
chmod +x "$PLUGIN_INSTALL_PATH"
sleep 0.5

echo "Executing the plugin"
# Uncomment the following line for debugging
export TF_LOG=DEBUG
terraform init && terraform validate
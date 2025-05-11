#!/bin/sh

# Variables for paths and plugin details
SRC_PATH="./cmd"
DEST_PATH="compiled_package"
PLUGIN_NAME="terraform-provider-ibm-verify"

# Clean old builds
echo "Cleaning old builds"
rm -rf "$DEST_PATH"/*

# Build artifacts for different platforms
echo "Building Mac - amd64 artifact"
GOOS=darwin GOARCH=amd64 go build -o "$DEST_PATH/${PLUGIN_NAME}_darwin_amd64" "$SRC_PATH"

echo "Building Mac - arm64 artifact"
GOOS=darwin GOARCH=arm64 go build -o "$DEST_PATH/${PLUGIN_NAME}_darwin_arm64" "$SRC_PATH"

echo "Building Linux - amd64 artifact"
GOOS=linux GOARCH=amd64 go build -o "$DEST_PATH/${PLUGIN_NAME}_linux_amd64" "$SRC_PATH"

# Change to the destination directory
cd "$DEST_PATH"

# Create ZIP files for each artifact
echo "Creating ZIP files for artifacts"
zip "${PLUGIN_NAME}_darwin_amd64.zip" "${PLUGIN_NAME}_darwin_amd64"
zip "${PLUGIN_NAME}_darwin_arm64.zip" "${PLUGIN_NAME}_darwin_arm64"
zip "${PLUGIN_NAME}_linux_amd64.zip" "${PLUGIN_NAME}_linux_amd64"

cd ..
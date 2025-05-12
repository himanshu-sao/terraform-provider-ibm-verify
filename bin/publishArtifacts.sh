#!/bin/bash

read -p "Enter the version (e.g., 0.0.1): " VERSION
mkdir -p artifacts
rm -rf artifacts/*

VERSION=$VERSION make

# Sign the SHA256SUMS file, please note that the default key is something specific to the user
cd . && gpg --default-key FF4107BA13B10B5DD64F515A7EA0CA77FFE0303A --detach-sign terraform-provider-ibmverify_${VERSION}_SHA256SUMS

# Move the files to the artifacts directory, excluding .exe files
find . -maxdepth 1 -type f -name "terraform-provider-ibmverify_${VERSION}_*" ! -name "*.exe" -exec mv -v {} artifacts/ \;
rm -rf "terraform-provider-ibmverify*.exe"
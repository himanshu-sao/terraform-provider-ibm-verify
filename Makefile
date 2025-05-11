# Variables
PLUGIN_NAME := terraform-provider-ibm-verify
PLUGIN_VERSION := 0.0.1
BUILD_DIR := ./compiled_package
CMD_DIR := ./cmd
BIN_DIR := ./bin
OS_ARCH := darwin_amd64
PLUGIN_INSTALL_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/local/ibmverify/$(PLUGIN_VERSION)/$(OS_ARCH)

# Default target
.PHONY: all
all: clean check build package

# Check code formatting and vet issues
.PHONY: check
check:
	@echo "Checking that code complies with gofmt requirements..."
	@go vet ./...

# Clean up old builds
.PHONY: clean
clean:
	@echo "Cleaning up old builds and artifacts..."
	rm -rf $(BUILD_DIR)/*
	rm -rf $(PLUGIN_INSTALL_DIR)
	rm -rf .terraform*

# Build the plugin
.PHONY: build
build: check
	@echo "Building the plugin..."
	@$(BIN_DIR)/buildPlugin.sh

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(BIN_DIR)/runTests.sh

# Package the plugin for multiple platforms
.PHONY: package
package:
	@echo "Packaging the plugin for distribution..."
	@$(BIN_DIR)/package.sh

# Install the plugin locally
.PHONY: install
install: build
	@echo "Installing the plugin locally..."
	mkdir -p $(PLUGIN_INSTALL_DIR)
	cp -v $(BUILD_DIR)/$(PLUGIN_NAME) $(PLUGIN_INSTALL_DIR)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)
	chmod +x $(PLUGIN_INSTALL_DIR)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)

# Validate the plugin with Terraform
.PHONY: validate
validate: install
	@echo "Validating the plugin with Terraform..."
	terraform init && terraform validate

# Full pipeline: clean, build, test, package, and validate
.PHONY: pipeline
pipeline: clean build test package validate
	@echo "Pipeline completed successfully!"
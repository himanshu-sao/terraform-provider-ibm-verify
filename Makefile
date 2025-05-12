# Variables
PLUGIN_NAME := terraform-provider-ibmverify
PLUGIN_VERSION := $(or $(VERSION), 0.0.1)
BUILD_DIR := ./compiled_package
ARTIFACTS_DIR := ./artifacts
OUTPUT_DIR := .
CMD_DIR := ./cmd
BIN_DIR := ./bin
OS_ARCHS := darwin_amd64 darwin_arm64 linux_amd64 linux_386 linux_arm linux_arm64 windows_amd64 windows_386
PLUGIN_INSTALL_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/local/ibmverify/$(PLUGIN_VERSION)/$(OS_ARCH)

# Default target
.PHONY: all
all: clean check build package checksums

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
	rm -rf $(OUTPUT_DIR)/$(PLUGIN_NAME)*
	rm -rf ${ARTIFACTS_DIR}/$(PLUGIN_NAME)*

# Build the plugin for all platforms
.PHONY: build
build: check
	@echo "Building the plugin for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for os_arch in $(OS_ARCHS); do \
		OS=$${os_arch%_*}; \
		ARCH=$${os_arch#*_}; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $(BUILD_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch $(CMD_DIR); \
		if [ "$$OS" = "windows" ]; then \
			mv $(BUILD_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch $(BUILD_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch.exe; \
		fi; \
	done

# Package the plugin for distribution
.PHONY: package
package: build
	@echo "Packaging the plugin for distribution..."
	@for os_arch in $(OS_ARCHS); do \
		OS=$${os_arch%_*}; \
		if [ "$$OS" = "windows" ]; then \
			cp $(BUILD_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch.exe $(OUTPUT_DIR)/; \
			zip -j $(OUTPUT_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch.zip $(OUTPUT_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch.exe; \
		else \
			cp $(BUILD_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch $(OUTPUT_DIR)/; \
			zip -j $(OUTPUT_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch.zip $(OUTPUT_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_$$os_arch; \
		fi; \
	done
	@echo "Packaging completed!"
	@echo "Plugin files are located in $(OUTPUT_DIR)"

# Generate checksums
.PHONY: checksums
checksums:
	@echo "Generating SHA256 checksums..."
	@cd $(OUTPUT_DIR) && shasum -a 256 $(PLUGIN_NAME)_$(PLUGIN_VERSION)_*.zip > $(PLUGIN_NAME)_$(PLUGIN_VERSION)_SHA256SUMS
	@echo "Checksums file created at $(OUTPUT_DIR)/$(PLUGIN_NAME)_$(PLUGIN_VERSION)_SHA256SUMS"
	@echo "Don't forget to sign the checksums file with:"
	@echo "cd $(OUTPUT_DIR) && gpg --detach-sign $(PLUGIN_NAME)_$(PLUGIN_VERSION)_SHA256SUMS"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(BIN_DIR)/runTests.sh

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
pipeline: clean build test package checksums validate
	@echo "Pipeline completed successfully!"
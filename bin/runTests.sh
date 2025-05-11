#!/bin/bash

# Variables for test paths
TEST_BASE="./pkg/test"
MODULE_BASE="./pkg/module"

# Run tests for oidc
echo "Running the tests for pkg/oidc"
go test "$TEST_BASE/oidc" -cover -coverpkg="$MODULE_BASE/oidc"

# Run tests for util
echo "Running the tests for pkg/util"
go test "$TEST_BASE/util" -cover -coverpkg="$MODULE_BASE/util"

# Run tests for branding
#echo "Running the tests for pkg/branding"
#go test "$TEST_BASE/branding" -cover -coverpkg="$MODULE_BASE/branding"
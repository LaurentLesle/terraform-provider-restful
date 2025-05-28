#!/bin/bash

# Build script for both Terraform and OpenTofu compatibility
# This script builds the provider for both registries

set -e

VERSION=${1:-"dev"}
PLATFORMS=${2:-"linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 windows_amd64"}

echo "Building restful provider v${VERSION} for platforms: ${PLATFORMS}"

# Create dist directory
mkdir -p dist

# Build for each platform
for platform in $PLATFORMS; do
    os=$(echo $platform | cut -d'_' -f1)
    arch=$(echo $platform | cut -d'_' -f2)
    
    echo "Building for $os/$arch..."
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    
    # Build binary
    binary_name="terraform-provider-restful_v${VERSION}"
    if [ "$os" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi
    
    go build -o "dist/${binary_name}" ./main.go
    
    # Create OpenTofu compatible binary (same binary, different name)
    opentofu_binary_name="opentofu-provider-restful_v${VERSION}"
    if [ "$os" = "windows" ]; then
        opentofu_binary_name="${opentofu_binary_name}.exe"
    fi
    
    cp "dist/${binary_name}" "dist/${opentofu_binary_name}"
    
    echo "✓ Built for $platform"
done

echo "Build complete! Binaries available in dist/"
echo ""
echo "Terraform binaries:"
ls -la dist/terraform-provider-restful_*
echo ""
echo "OpenTofu binaries:"
ls -la dist/opentofu-provider-restful_*

#!/bin/bash

# Script to update all imports from terraform-provider-restful to opentofu-provider-restful

echo "Updating import paths..."

# Find all Go files and update the module import paths
find . -name "*.go" -type f -exec sed -i 's|github.com/magodo/terraform-provider-restful|github.com/magodo/opentofu-provider-restful|g' {} \;

echo "Import paths updated successfully!"

# Test that the changes work
echo "Running go mod tidy to verify dependencies..."
go mod tidy

if [ $? -eq 0 ]; then
    echo "✅ Go mod tidy successful - dependencies are correct"
else
    echo "❌ Go mod tidy failed - there might be import issues"
    exit 1
fi

echo "Running go build to verify compilation..."
go build -o /tmp/test-build ./main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful - code compiles correctly"
    rm -f /tmp/test-build
else
    echo "❌ Build failed - there might be compilation issues"
    exit 1
fi

echo "🎉 OpenTofu conversion completed successfully!"

#!/bin/bash

# Update all import statements from magodo to aztfmod namespace
find . -name "*.go" -type f -exec sed -i 's|github.com/magodo/opentofu-provider-restful|github.com/aztfmod/opentofu-provider-restful|g' {} \;

echo "Updated all Go files with new aztfmod namespace"

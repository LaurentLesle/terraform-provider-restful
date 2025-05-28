#!/bin/bash

# Update all import statements from aztfmod to aztfmod namespace
find . -name "*.go" -type f -exec sed -i 's|github.com/aztfmod/opentofu-provider-restful|github.com/aztfmod/opentofu-provider-restful|g' {} \;

echo "Updated all Go files with new aztfmod namespace"

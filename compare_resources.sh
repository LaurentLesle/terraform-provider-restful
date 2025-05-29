#!/bin/bash

# Script to compare resource_definitions.json against RESOURCE_TYPES.md
# and identify missing/incorrect slugs

echo "=== COMPARISON: resource_definitions.json vs RESOURCE_TYPES.md ==="
echo ""

# Extract resource names from RESOURCE_TYPES.md (skipping header)
echo "Extracting resources from RESOURCE_TYPES.md..."
reference_resources=$(grep "^| azurerm_" /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | awk -F'|' '{print $2}' | sed 's/^ *//g' | sed 's/ *$//g')

# Extract resource names from JSON (get keys)
echo "Extracting resources from resource_definitions.json..."
json_resources=$(jq -r 'keys[]' /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json | grep "^azurerm_")

echo ""
echo "=== RESOURCES IN REFERENCE BUT MISSING FROM JSON ==="
echo ""

missing_count=0
for resource in $reference_resources; do
    if ! echo "$json_resources" | grep -q "^$resource$"; then
        echo "MISSING: $resource"
        # Get the slug from RESOURCE_TYPES.md
        slug=$(grep "| $resource " /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | awk -F'|' '{print $3}' | sed 's/^ *`//g' | sed 's/`.*$//g')
        echo "  Suggested slug: $slug"
        echo ""
        ((missing_count++))
    fi
done

echo ""
echo "=== RESOURCES IN JSON BUT NOT IN REFERENCE ==="
echo ""

extra_count=0
for resource in $json_resources; do
    if ! echo "$reference_resources" | grep -q "^$resource$"; then
        echo "EXTRA: $resource"
        # Get the slug from JSON
        slug=$(jq -r ".\"$resource\".slug" /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json)
        echo "  Current slug: $slug"
        echo ""
        ((extra_count++))
    fi
done

echo ""
echo "=== SLUG MISMATCHES (same resource, different slug) ==="
echo ""

mismatch_count=0
for resource in $json_resources; do
    if echo "$reference_resources" | grep -q "^$resource$"; then
        # Both have this resource, compare slugs
        ref_slug=$(grep "| $resource " /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | awk -F'|' '{print $3}' | sed 's/^ *`//g' | sed 's/`.*$//g')
        json_slug=$(jq -r ".\"$resource\".slug" /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json)
        
        if [ "$ref_slug" != "$json_slug" ]; then
            echo "MISMATCH: $resource"
            echo "  Reference slug: $ref_slug"
            echo "  JSON slug: $json_slug"
            echo ""
            ((mismatch_count++))
        fi
    fi
done

echo ""
echo "=== SUMMARY ==="
echo "Missing from JSON: $missing_count"
echo "Extra in JSON: $extra_count"
echo "Slug mismatches: $mismatch_count"
echo ""

# Let's also count total resources
total_ref=$(echo "$reference_resources" | wc -l)
total_json=$(echo "$json_resources" | wc -l)
echo "Total in reference: $total_ref"
echo "Total in JSON: $total_json"

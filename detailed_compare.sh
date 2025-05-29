#!/bin/bash

echo "=== DETAILED RESOURCE COMPARISON ==="
echo ""

# Get unique resource names from RESOURCE_TYPES.md
echo "Getting unique resources from RESOURCE_TYPES.md..."
reference_resources=$(grep "^| azurerm_" /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | awk -F'|' '{print $2}' | sed 's/^ *//g' | sed 's/ *$//g' | sort | uniq)

# Get resource names from JSON
echo "Getting resources from resource_definitions.json..."
json_resources=$(jq -r 'keys[]' /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json | grep "^azurerm_" | sort)

echo ""
echo "=== MISSING FROM JSON (found in RESOURCE_TYPES.md but not in JSON) ==="
echo ""

missing_count=0
for resource in $reference_resources; do
    if ! echo "$json_resources" | grep -q "^$resource$"; then
        echo "MISSING: $resource"
        # Get the slug from RESOURCE_TYPES.md (first occurrence)
        slug=$(grep "| $resource " /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | head -1 | awk -F'|' '{print $3}' | sed 's/^ *`//g' | sed 's/`.*$//g')
        echo "  Suggested slug: $slug"
        echo ""
        ((missing_count++))
    fi
done

echo ""
echo "=== MISSING FROM REFERENCE (found in JSON but not in RESOURCE_TYPES.md) ==="
echo ""

extra_count=0
for resource in $json_resources; do
    if ! echo "$reference_resources" | grep -q "^$resource$"; then
        echo "NOT IN REFERENCE: $resource"
        # Get the slug from JSON
        slug=$(jq -r ".\"$resource\".slug" /workspaces/terraform-provider-restful/internal/provider/resource_definitions.json 2>/dev/null)
        echo "  Current slug: $slug"
        echo ""
        ((extra_count++))
    fi
done

echo ""
echo "=== DUPLICATES IN REFERENCE FILE ==="
echo ""

# Check for duplicates in reference
echo "Checking for duplicate resource names in RESOURCE_TYPES.md..."
duplicate_resources=$(grep "^| azurerm_" /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | awk -F'|' '{print $2}' | sed 's/^ *//g' | sed 's/ *$//g' | sort | uniq -d)

if [ -n "$duplicate_resources" ]; then
    for dup_resource in $duplicate_resources; do
        echo "DUPLICATE: $dup_resource"
        grep "| $dup_resource " /workspaces/terraform-provider-restful/RESOURCE_TYPES.md | while IFS= read -r line; do
            slug=$(echo "$line" | awk -F'|' '{print $3}' | sed 's/^ *`//g' | sed 's/`.*$//g')
            echo "  Variant slug: $slug"
        done
        echo ""
    done
else
    echo "No duplicates found in reference file."
fi

echo ""
echo "=== SUMMARY ==="
echo "Missing from JSON: $missing_count"
echo "Extra in JSON (not in reference): $extra_count"

# Count totals
total_ref_unique=$(echo "$reference_resources" | wc -l)
total_json=$(echo "$json_resources" | wc -l)
echo "Total unique in reference: $total_ref_unique"
echo "Total in JSON: $total_json"

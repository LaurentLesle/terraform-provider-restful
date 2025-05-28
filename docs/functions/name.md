---
page_title: "name Function - terraform-provider-restful"
subcategory: "Functions"
description: |-
  Generate standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions.
---

# name Function

The `name` function generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This implements functionality similar to `azurecaf_name` from the terraform-provider-azurecaf.

## Syntax

```hcl
provider::restful::name(name, resource_type)
provider::restful::name(name, resource_type, options...)
```

## Arguments

1. `name` (String) - The base name for the resource.
2. `resource_type` (String) - The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account').
3. `options` (List of Strings, Optional) - Optional parameters for customizing name generation.

## Supported Resource Types

The function currently supports the following Azure resource types:

- `azurerm_resource_group` - Resource Group (slug: "rg")
- `azurerm_storage_account` - Storage Account (slug: "st")  
- `azurerm_key_vault` - Key Vault (slug: "kv")
- `azurerm_container_registry` - Container Registry (slug: "cr")
- `azurerm_virtual_network` - Virtual Network (slug: "vnet")

## Examples

### Basic Usage

```hcl
locals {
  # Generate a storage account name: "stmyapp"
  storage_name = provider::restful::name("myapp", "azurerm_storage_account")
  
  # Generate a resource group name: "rg-myproject"
  rg_name = provider::restful::name("myproject", "azurerm_resource_group")
  
  # Generate a key vault name: "kv-secrets"
  kv_name = provider::restful::name("secrets", "azurerm_key_vault")
}
```

### With Configuration Options

```hcl
# Note: Options parameter implementation is planned for future versions
# This will support prefixes, suffixes, random_length, etc.
locals {
  # Future syntax example:
  # complex_name = provider::restful::name("myapp", "azurerm_storage_account", [
  #   "prefixes=[\"dev\", \"west\"]",
  #   "random_length=4",
  #   "separator=-"
  # ])
}
```

## Resource Naming Rules

Each resource type has specific naming rules that are automatically applied:

### Storage Account (`azurerm_storage_account`)
- Length: 3-24 characters
- Characters: lowercase letters and numbers only
- Pattern: `^[a-z0-9]{3,24}$`

### Resource Group (`azurerm_resource_group`)
- Length: 1-90 characters  
- Characters: letters, numbers, hyphens, periods, underscores, parentheses
- Cannot end with period

### Key Vault (`azurerm_key_vault`)
- Length: 3-24 characters
- Characters: letters, numbers, hyphens
- Must start with letter, end with letter or number

### Container Registry (`azurerm_container_registry`)
- Length: 5-50 characters
- Characters: letters and numbers only
- Pattern: `^[a-zA-Z0-9]{5,50}$`

### Virtual Network (`azurerm_virtual_network`)
- Length: 2-64 characters
- Characters: letters, numbers, hyphens, periods, underscores
- Must start with letter or number, end with letter, number, or underscore

## Implementation Details

The function automatically:

1. **Adds resource slug**: Prepends the appropriate slug (e.g., "st" for storage accounts)
2. **Applies cleaning**: Removes invalid characters based on resource type rules
3. **Enforces case**: Converts to lowercase for resources that require it
4. **Validates length**: Ensures the generated name meets minimum and maximum length requirements
5. **Validates pattern**: Checks the final name against the resource type's validation regex

## Return Value

Returns a string containing the generated resource name that complies with Azure naming conventions for the specified resource type.

## Notes

- The function generates deterministic names - the same inputs will always produce the same output
- Names are automatically validated against Azure naming conventions
- Invalid characters are automatically removed based on the resource type
- Future versions will support advanced options like prefixes, suffixes, and random components

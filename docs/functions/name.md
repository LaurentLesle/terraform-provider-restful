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
provider::restful::name({
  name          = "base-name"
  resource_type = "azure_resource_type"
  # ... optional parameters
})
```

## Arguments

The function takes a single configuration object with the following attributes:

### Required Attributes

- `name` (String) - The base name for the resource
- `resource_type` (String) - The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account')

### Optional Attributes

- `prefix` / `prefixes` (String or Array) - Single prefix string or array of prefixes to prepend
- `suffix` / `suffixes` (String or Array) - Single suffix string or array of suffixes to append  
- `separator` (String) - Character(s) used to separate name components (default: "-")
- `clean_input` (Boolean) - Whether to clean input strings according to resource rules (default: true)
- `use_slug` (Boolean) - Whether to include the resource type slug in the name (default: true)
- `passthrough` (Boolean) - Whether to skip processing and return the name as-is (default: false)
- `random_length` (Number) - Length of random suffix to append (default: 0, no random suffix)

## Supported Resource Types

The function currently supports the following Azure resource types:

- `azurerm_resource_group` - Resource Group (slug: "rg")
- `azurerm_storage_account` - Storage Account (slug: "st")  
- `azurerm_key_vault` - Key Vault (slug: "kv")
- `azurerm_container_registry` - Container Registry (slug: "cr")
- `azurerm_virtual_network` - Virtual Network (slug: "vnet")

## Examples

### Simple Usage

```hcl
# Basic resource group name
output "rg_name" {
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
  })
  # Result: "rg-myproject"
}

# Basic storage account name
output "storage_name" {
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
  })
  # Result: "stmyapp"
}

# Key vault with prefix
output "kv_name" {
  value = provider::restful::name({
    name          = "secrets"
    resource_type = "azurerm_key_vault"
    prefix        = "dev"
  })
  # Result: "dev-kv-secrets"
}
```

### Advanced Usage

```hcl
# Storage account with multiple suffixes and no separator
output "storage_complex" {
  value = provider::restful::name({
    name          = "data"
    resource_type = "azurerm_storage_account"
    prefixes      = ["prod", "team1"]
    suffixes      = ["001", "cache"]
    separator     = ""
  })
  # Result: "prodteam1stdatacache001"
}

# Resource group with random suffix
output "rg_with_random" {
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
    random_length = 3
  })
  # Result: "dev-rg-webapp-001-r96" (random suffix varies)
}

# Passthrough mode (skip all processing)
output "custom_name" {
  value = provider::restful::name({
    name          = "MyCustomName"
    resource_type = "azurerm_resource_group"
    passthrough   = true
  })
  # Result: "MyCustomName"
}

# Custom separator
output "underscore_separated" {
  value = provider::restful::name({
    name          = "api"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
    separator     = "_"
  })
  # Result: "dev_rg_api_001"
}
```

### Error Handling

```hcl
# This will fail with unauthorized attribute error
output "invalid_config" {
  value = provider::restful::name({
    name          = "test"
    resource_type = "azurerm_resource_group"
    invalid_attr  = "value"  # This will cause an error
  })
  # Error: Unauthorized attributes found: 'invalid_attr' (not supported)
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

1. **Validates input attributes**: Checks that all provided attributes are authorized and returns clear error messages for unauthorized attributes
2. **Adds resource slug**: Prepends the appropriate slug (e.g., "st" for storage accounts) when `use_slug` is true
3. **Applies cleaning**: Removes invalid characters based on resource type rules when `clean_input` is true
4. **Enforces case**: Converts to lowercase for resources that require it
5. **Validates length**: Ensures the generated name meets minimum and maximum length requirements
6. **Validates pattern**: Checks the final name against the resource type's validation regex
7. **Supports flexible parameters**: Accepts both singular (`prefix`, `suffix`) and plural (`prefixes`, `suffixes`) forms
8. **Generates random suffixes**: Adds deterministic random strings when `random_length` is specified

## Attribute Validation

The function validates all input attributes and will return an error for any unauthorized attributes. The error message will explicitly list the unauthorized attributes with "(not supported)" comments and provide the complete list of supported attributes.

**Supported attributes:**
- `name`, `resource_type` (required)
- `prefix`, `prefixes`, `suffix`, `suffixes`
- `separator`, `clean_input`, `use_slug`, `passthrough`, `random_length`

## Return Value

Returns a string containing the generated resource name that complies with Azure naming conventions for the specified resource type.

## Notes

- The function generates deterministic names - the same inputs will always produce the same output
- Random suffixes use a deterministic seed, so they're predictable in tests but unique in practice
- Names are automatically validated against Azure naming conventions
- Invalid characters are automatically removed based on the resource type
- Both singular and plural parameter forms are supported for flexibility (`prefix`/`prefixes`, `suffix`/suffixes`)
- The function provides comprehensive error messages for troubleshooting

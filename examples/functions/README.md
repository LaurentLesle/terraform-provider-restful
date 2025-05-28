# Name Function Examples

This directory contains examples demonstrating the `provider::restful::name` function for generating standardized Azure resource names.

## Files

- **`name.tf`** - Main examples showing both simple and advanced usage
- **`comprehensive_examples.tf`** - Additional examples and error handling demonstrations

## Quick Start

### Simple Usage
```hcl
# Basic resource group
output "rg_name" {
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
  })
  # Result: "rg-myproject"
}

# Basic storage account
output "storage_name" {
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    separator     = ""
  })
  # Result: "stmyapp"
}
```

### Advanced Usage
```hcl
# Complex configuration
output "complex_name" {
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_resource_group"
    prefix        = "prod"
    suffix        = "001"
    random_length = 3
  })
  # Result: "prod-rg-webapp-001-r96"
}
```

## Supported Attributes

**Required:**
- `name` - Base name for the resource
- `resource_type` - Azure resource type

**Optional:**
- `prefix`/`prefixes` - Single or multiple prefixes
- `suffix`/`suffixes` - Single or multiple suffixes
- `separator` - Component separator (default: "-")
- `clean_input` - Clean input strings (default: true)
- `use_slug` - Include resource type slug (default: true)
- `passthrough` - Skip processing (default: false)
- `random_length` - Add random suffix (default: 0)

## Error Handling

The function validates all attributes and provides clear error messages for unauthorized attributes:

```
Error: Unauthorized attributes found: 'invalid_attr' (not supported).
Supported attributes are: name, resource_type, prefix/prefixes, suffix/suffixes, separator, clean_input, use_slug, passthrough, random_length
```

## Running Examples

```bash
# Initialize and plan
cd examples/functions
export TF_CLI_CONFIG_FILE=$(pwd)/.tofurc
tofu plan

# Apply changes
tofu apply -auto-approve
```

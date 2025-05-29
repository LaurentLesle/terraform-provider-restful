# Comprehensive examples showing the name function capabilities
# This file demonstrates both valid usage and error cases

# === VALID EXAMPLES ===

# All supported attributes example
output "all_features_demo" {
  description = "Demonstrates all supported attributes"
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_resource_group" 
    prefix        = "prod"
    suffix        = "001"
    separator     = "-"
    clean_input   = true
    use_slug      = true
    passthrough   = false
    random_length = 3
  })
  # Result: "prod-rg-webapp-001-r96"
}

# Singular vs plural forms
output "singular_form" {
  description = "Using singular prefix/suffix"
  value = provider::restful::name({
    name          = "api"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
  })
  # Result: "dev-rg-api-001"
}

output "plural_form" {
  description = "Using plural prefixes/suffixes"
  value = provider::restful::name({
    name          = "api"
    resource_type = "azurerm_resource_group"
    prefixes      = ["dev", "west"]
    suffixes      = ["001", "cache"]
  })
  # Result: "dev-west-rg-api-001-cache"
}

# Resource-specific examples
output "storage_with_random" {
  description = "Storage account with random suffix"
  value = provider::restful::name({
    name          = "logs"
    resource_type = "azurerm_storage_account"
    prefix        = "prod"
    random_length = 4
    separator     = ""
  })
  # Result: "prodstlogsr96x" (4-char random suffix)
}

# === ERROR DEMONSTRATION (commented out to avoid plan failures) ===

# Uncomment any of these to see error messages:

# output "unauthorized_attribute_error" {
#   description = "This will fail - bla is not a supported attribute"
#   value = provider::restful::name({
#     name          = "test"
#     resource_type = "azurerm_resource_group"
#     bla           = "invalid"  # This causes error
#   })
#   # Error: Unauthorized attributes found: 'bla' (not supported). 
#   # Supported attributes are: name, resource_type, prefix/prefixes, suffix/suffixes, separator, clean_input, use_slug, passthrough, random_length
# }

# output "multiple_unauthorized_error" {
#   description = "Multiple unauthorized attributes"
#   value = provider::restful::name({
#     name          = "test"
#     resource_type = "azurerm_resource_group"
#     invalid_one   = "value1"
#     invalid_two   = "value2"
#     also_invalid  = "value3"
#   })
#   # Error: Unauthorized attributes found: 'invalid_one' (not supported), 'invalid_two' (not supported), 'also_invalid' (not supported).
# }

# output "missing_required_error" {
#   description = "Missing required attribute"
#   value = provider::restful::name({
#     # name missing - will cause error
#     resource_type = "azurerm_resource_group"
#   })
#   # Error: The 'name' attribute is required in settings
# }

# output "invalid_resource_type_error" {
#   description = "Invalid resource type"
#   value = provider::restful::name({
#     name          = "test"
#     resource_type = "invalid_resource_type"
#   })
#   # Error: Error generating name: unsupported resource type: invalid_resource_type
# }

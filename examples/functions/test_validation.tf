# Test validation of authorized and unauthorized attributes

# This should work - random_length is authorized
output "valid_with_random_length" {
  description = "Valid usage with random_length"
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    random_length = 3
    separator     = ""
  })
}

# This should fail - bla is not authorized
# output "invalid_with_bla" {
#   description = "Invalid usage with unauthorized attribute"
#   value = provider::restful::name({
#     name          = "myapp"
#     resource_type = "azurerm_storage_account"
#     bla           = "xxx"
#   })
# }

# Comprehensive test with mandatory resource_type

output "test_mandatory_fields" {
  description = "Test with only mandatory fields"
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_app_service"
  })
}

output "test_storage_with_arrays" {
  description = "Test storage account with prefixes and suffixes"
  value = provider::restful::name({
    name          = "data"
    resource_type = "azurerm_storage_account"
    prefixes      = ["prod", "analytics"]
    suffixes      = ["001", "cache"]
    separator     = ""
  })
}

output "test_keyvault_singular" {
  description = "Test key vault with singular prefix/suffix"
  value = provider::restful::name({
    name          = "secrets"
    resource_type = "azurerm_key_vault"
    prefix        = "shared"
    suffix        = "main"
  })
}

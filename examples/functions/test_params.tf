# Test both singular and plural parameter forms

output "test_singular" {
  description = "Test with singular prefix/suffix"
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_resource_group"
    prefix        = "prod"
    suffix        = "001"
  })
}

output "test_plural" {
  description = "Test with plural prefixes/suffixes"
  value = provider::restful::name({
    name          = "webapp"
    resource_type = "azurerm_resource_group"
    prefixes      = ["prod", "team1"]
    suffixes      = ["001", "main"]
  })
}

output "test_mixed_arrays" {
  description = "Test with mixed array usage"
  value = provider::restful::name({
    name          = "api"
    resource_type = "azurerm_storage_account"
    prefix        = "dev"
    suffixes      = ["001", "cache"]
    separator     = ""
  })
}

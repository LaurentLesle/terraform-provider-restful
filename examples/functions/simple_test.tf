# Simple test
output "simple_name" {
  value = provider::restful::name({
    name          = "test"
    resource_type = "azurerm_resource_group"
  })
}

# Resource group with custom settings
output "rg_name" {
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_resource_group"
    prefixes      = ["prod", "team1"]
    suffixes      = ["main"]
  })
}

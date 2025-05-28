# Example: Using the name function

terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://api.example.com"
}

# === SIMPLE EXAMPLES ===

# Basic resource group name
output "simple_rg_name" {
  description = "Simple resource group name"
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
  })
  # Result: "rg-myproject"
}

# Basic storage account name  
output "simple_storage_name" {
  description = "Simple storage account name"
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    separator     = ""
  })
  # Result: "stmyapp"
}

# === ADVANCED EXAMPLES ===

# Storage account with complex configuration
output "storage_account_name" {
  description = "Generated storage account name with suffixes"
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    suffixes      = ["001"]
    separator     = ""
    use_slug      = true
  })
  # Result: "stmyapp001"
}

# Resource group with prefix, suffix, and random component
output "resource_group_name" {
  description = "Generated resource group name with random suffix"
  value = provider::restful::name({
      name          = "myproject"
      resource_type = "azurerm_resource_group"
      prefix        = "dev"
      suffix        = "001"
      random_length = 3
      passthrough   = false
      use_slug      = true
  })
  # Result: "dev-rg-myproject-001-r96"
}

# Multiple prefixes and suffixes
output "complex_storage_name" {
  description = "Complex storage account with multiple prefixes and suffixes"
  value = provider::restful::name({
    name          = "data"
    resource_type = "azurerm_storage_account"
    prefixes      = ["prod", "team1"]
    suffixes      = ["cache", "001"]
    separator     = ""
  })
  # Result: "prodteam1stdatacache001"
}

# Custom separator example
output "underscore_separated_name" {
  description = "Resource group with custom separator"
  value = provider::restful::name({
    name          = "api"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
    separator     = "_"
  })
  # Result: "dev_rg_api_001"
}

# Passthrough mode example
output "passthrough_name" {
  description = "Name with passthrough mode (no processing)"
  value = provider::restful::name({
    name          = "MyCustomName"
    resource_type = "azurerm_resource_group"
    passthrough   = true
  })
  # Result: "MyCustomName"
}


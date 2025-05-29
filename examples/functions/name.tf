# Example: Using the name function for Azure resource naming

terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://api.example.com"  # Required but not used for functions
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
    name          = "webapp"
    resource_type = "azurerm_storage_account"
    prefixes      = ["dev", "east"]
    suffixes      = ["01"]
    separator     = ""
  })
  # Result: "deveaststwebapp01"
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
  # Result: "dev-rg-myproject-001-abc"
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

# Virtual machine naming
output "vm_name" {
  description = "Virtual machine name with prefix"
  value = provider::restful::name({
    name          = "worker"
    resource_type = "azurerm_linux_virtual_machine"
    prefix        = "dev"
    clean_input   = true
  })
  # Result: "dev-vm-worker"
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

# Random suffix example
output "random_storage_name" {
  description = "Storage account with random suffix"
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    random_length = 3
    separator     = ""
  })
  # Result: "stmyappbsz" (random suffix will vary with seed)
}

# Azure AD B2C Directory example
output "aad_b2c_name" {
  description = "Azure AD B2C Directory name"
  value = provider::restful::name({
    name          = "aad"
    resource_type = "azurerm_aadb2c_directory"
    random_length = 3
  })
  # Result: "aad-bsz" (with random suffix)
}

# === TESTING DIFFERENT CONFIGURATIONS ===

# Test with slug disabled
output "no_slug_name" {
  description = "Resource group name without slug"
  value = provider::restful::name({
    name          = "test"
    resource_type = "azurerm_resource_group"
    use_slug      = false
  })
  # Result: "test"
}

# Test with clean input disabled
output "no_clean_name" {
  description = "Name with valid characters (clean_input disabled)"
  value = provider::restful::name({
    name          = "my-app-123"
    resource_type = "azurerm_resource_group"
    clean_input   = false
  })
  # Result: "rg-my-app-123" (valid, no cleaning needed)
}


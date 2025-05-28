# Example using the restful_name data source to generate Azure resource names

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

# Generate a resource group name
data "restful_name" "rg_example" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev", "eastus"]
  suffixes      = ["001"]
  random_length = 3
  random_seed   = 12345
  clean_input   = true
}

# Generate a storage account name
data "restful_name" "storage_example" {
  name          = "myappstorage"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  random_length = 5
  random_seed   = 67890
  clean_input   = true
}

# Generate a key vault name
data "restful_name" "kv_example" {
  name          = "myapp"
  resource_type = "azurerm_key_vault"
  prefixes      = ["dev"]
  suffixes      = ["secrets"]
  clean_input   = true
}

# Generate a virtual network name
data "restful_name" "vnet_example" {
  name          = "myapp"
  resource_type = "azurerm_virtual_network"
  prefixes      = ["dev", "eastus"]
  clean_input   = true
}

# Generate a container registry name
data "restful_name" "acr_example" {
  name          = "myappregistry"
  resource_type = "azurerm_container_registry"
  prefixes      = ["dev"]
  random_length = 4
  clean_input   = true
}

# Example using passthrough mode (validation only)
data "restful_name" "passthrough_example" {
  name        = "my-existing-resource"
  resource_type = "azurerm_resource_group"
  passthrough = true
  clean_input = true
}

# Outputs to show the generated names
output "resource_group_name" {
  value       = data.restful_name.rg_example.result
  description = "Generated resource group name"
}

output "storage_account_name" {
  value       = data.restful_name.storage_example.result
  description = "Generated storage account name"
}

output "key_vault_name" {
  value       = data.restful_name.kv_example.result
  description = "Generated key vault name"
}

output "virtual_network_name" {
  value       = data.restful_name.vnet_example.result
  description = "Generated virtual network name"
}

output "container_registry_name" {
  value       = data.restful_name.acr_example.result
  description = "Generated container registry name"
}

output "passthrough_name" {
  value       = data.restful_name.passthrough_example.result
  description = "Validated name using passthrough mode"
}

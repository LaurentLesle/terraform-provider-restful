# Example: Using the restful_name data source for Azure resource naming
# Note: This file should be used together with name.tf or copy the terraform/provider blocks

# === BASIC DATA SOURCE EXAMPLES ===

# Basic resource group naming
data "restful_name" "rg_basic" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
}

# Basic storage account naming with custom settings
data "restful_name" "storage_basic" {
  name          = "mydata"
  resource_type = "azurerm_storage_account"
  separator     = ""
  clean_input   = true
}

# === ADVANCED DATA SOURCE EXAMPLES ===

# Storage account with prefixes, suffixes, and random component
data "restful_name" "storage_advanced" {
  name          = "mydata"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod", "eastus"]
  suffixes      = ["logs"]
  random_length = 4
  random_seed   = 12345  # Fixed seed for consistent results
  separator     = ""
}

# Virtual machine with prefix and clean input
data "restful_name" "vm_data_source" {
  name          = "worker-node"
  resource_type = "azurerm_linux_virtual_machine"
  prefixes      = ["dev"]
  clean_input   = true
}

# Resource group with custom separator
data "restful_name" "rg_custom_sep" {
  name          = "api-service"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  separator     = "_"
}

# Storage account with slug disabled
data "restful_name" "storage_no_slug" {
  name          = "customname"
  resource_type = "azurerm_storage_account"
  use_slug      = false
  separator     = ""
}

# Passthrough mode example
data "restful_name" "passthrough_data_source" {
  name          = "MyExistingResourceName"
  resource_type = "azurerm_resource_group"
  passthrough   = true
  clean_input   = false
}

# Azure Cosmos DB with multiple prefixes and suffixes
data "restful_name" "cosmos_db" {
  name          = "documents"
  resource_type = "azurerm_cosmosdb_account"
  prefixes      = ["prod", "global"]
  suffixes      = ["primary", "001"]
  random_length = 3
  separator     = "-"
}

# === OUTPUT THE RESULTS ===

output "data_source_rg_basic_name" {
  description = "Basic resource group name from data source"
  value       = data.restful_name.rg_basic.result
  # Expected: "rg-myapp"
}

output "data_source_storage_basic_name" {
  description = "Basic storage account name from data source"
  value       = data.restful_name.storage_basic.result
  # Expected: "stmydata"
}

output "data_source_storage_advanced_name" {
  description = "Advanced storage account name with random suffix"
  value       = data.restful_name.storage_advanced.result
  # Expected: "prodeastusstmlogsdataXXXX" (XXXX = random with seed 12345)
}

output "data_source_vm_name" {
  description = "Virtual machine name from data source"
  value       = data.restful_name.vm_data_source.result
  # Expected: "dev-vm-workernode" (cleaned input)
}

output "data_source_custom_separator_name" {
  description = "Resource group with custom separator from data source"
  value       = data.restful_name.rg_custom_sep.result
  # Expected: "prod_rg_api_service_001"
}

output "data_source_no_slug_storage_name" {
  description = "Storage account without slug from data source"
  value       = data.restful_name.storage_no_slug.result
  # Expected: "customname"
}

output "data_source_passthrough_name" {
  description = "Passthrough mode result from data source"
  value       = data.restful_name.passthrough_data_source.result
  # Expected: "MyExistingResourceName"
}

output "data_source_cosmos_db_name" {
  description = "Cosmos DB account name from data source"
  value       = data.restful_name.cosmos_db.result
  # Expected: "prod-global-cosmos-documents-primary-001-XXX"
}

# === DEMONSTRATION OF ID OUTPUT ===

output "data_source_storage_resource_id" {
  description = "Data source ID (same as result)"
  value       = data.restful_name.storage_advanced.id
}

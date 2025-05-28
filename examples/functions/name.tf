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

# Output the generated names using the function directly
output "storage_account_name" {
  description = "Generated storage account name"
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    suffixes      = ["001"]
    separator     = ""
    use_slug      = true
  })
}

output "resource_group_name" {
  description = "Generated resource group name"
  value = provider::restful::name({
      name          = "myproject"
      resource_type = "azurerm_resource_group"
      prefix        = "dev"
      suffix        = "001"
      random_length = 3
      passthrough   = false
      use_slug      = true
  })
}


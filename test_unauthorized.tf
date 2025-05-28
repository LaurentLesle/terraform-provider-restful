# Test unauthorized attribute validation

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

# This should fail - bla is not authorized
output "test_unauthorized" {
  value = provider::restful::name({
    name          = "myapp"
    resource_type = "azurerm_storage_account"
    bla           = "xxx"
  })
}

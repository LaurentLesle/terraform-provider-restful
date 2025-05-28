terraform {
  required_providers {
    restful = {
      source  = "aztfmod/restful"
      version = "~> 1.0"
    }
  }
}

provider "restful" {
  base_url = "https://api.example.com"
}

# This should trigger an error due to unauthorized attributes
output "test_with_invalid_attr" {
  description = "This should fail with unauthorized attribute error"
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
    random_length = 3        # This should cause an error (not authorized)
    bla          = "xxx"     # This should cause an error (not authorized)
    passthrough   = false
    use_slug      = true
  })
}

# This should work fine (only authorized attributes)
output "test_valid" {
  description = "This should work fine"
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
    prefix        = "dev"
    suffix        = "001"
    passthrough   = false
    use_slug      = true
  })
}

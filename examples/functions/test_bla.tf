output "test_valid_with_random_length" {
  value = provider::restful::name({
    name = "test"
    resource_type = "azurerm_resource_group"
    random_length = 3
  })
}

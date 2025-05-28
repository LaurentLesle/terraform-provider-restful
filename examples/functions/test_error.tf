# Test missing resource_type
output "test_missing_resource_type" {
  value = provider::restful::name({
    name = "test"
  })
}

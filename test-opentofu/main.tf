terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://httpbin.org"
}

# Test data source
data "restful_resource" "test" {
  id = "/get"
}

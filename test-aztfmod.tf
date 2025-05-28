terraform {
  required_providers {
    restful = {
      source = "registry.opentofu.org/aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://example.com/api"
}

# OpenTofu Compatibility Guide

This provider is fully compatible with OpenTofu, the open-source fork of Terraform. You can use it exactly the same way as with Terraform.

## Installation

### Using OpenTofu CLI

The provider is available on the OpenTofu provider registry:

```hcl
terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
      version = "~> 1.0"
    }
  }
}
```

### Registry Configuration

For OpenTofu, the provider is available at:
- OpenTofu Registry: `registry.opentofu.org/aztfmod/restful`
- Terraform Registry: `registry.terraform.io/aztfmod/restful` (also works with OpenTofu)

## Migration from Terraform

If you're migrating an existing Terraform configuration to OpenTofu, no changes are required to your provider configuration. All resources, data sources, and configuration options work identically.

### Before (Terraform)
```hcl
terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://api.example.com"
  security = {
    http = {
      token = {
        token = var.api_token
      }
    }
  }
}
```

### After (OpenTofu)
```hcl
# Exactly the same configuration works with OpenTofu
terraform {
  required_providers {
    restful = {
      source = "aztfmod/restful"
    }
  }
}

provider "restful" {
  base_url = "https://api.example.com"
  security = {
    http = {
      token = {
        token = var.api_token
      }
    }
  }
}
```

## Features Support

All features of the provider are fully supported in OpenTofu:

- ✅ HTTP Authentication (Basic, Token)
- ✅ API Key Authentication  
- ✅ OAuth2 (Client Credentials, Password, Refresh Token)
- ✅ Custom CRUD methods and paths
- ✅ Precheck conditions
- ✅ Polling asynchronous operations
- ✅ Partial body tracking
- ✅ Operation resources
- ✅ Ephemeral resources
- ✅ Write-only attributes
- ✅ All import capabilities

## Examples

See the `examples/` directory.

## State Migration

If you're migrating an existing Terraform state to OpenTofu:

1. Copy your `.terraform/` directory and state files
2. Run `tofu init` to initialize OpenTofu
3. Continue using `tofu` commands instead of `terraform`

Your existing state and configurations will work without modification.

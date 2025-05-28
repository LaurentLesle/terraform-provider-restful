# Terraform/OpenTofu Provider Restful

This is a general Terraform and OpenTofu provider that aims to work for any platform as long as it exposes a RESTful API.

The document of this provider is available on:
- [Terraform Provider Registry](https://registry.terraform.io/providers/aztfmod/restful/latest/docs)
- [OpenTofu Provider Registry](https://registry.opentofu.org/providers/aztfmod/restful/latest/docs)

## Features

- Different authentication choices: HTTP auth (basic, token), API Key auth and OAuth2 (client credential, password credential, refresh token).
- Customized CRUD methods and paths
- Support precheck conditions
- Support polling asynchronous operations
- Partial `body` tracking: only the specified properties of the resource in the `body` attribute is tracked for diffs
- `restful_operation` resource that supports arbitrary Restful API call (e.g. `POST`) on create/update
- Ephemeral resource `restful_resource`
- [Write-only attributes](https://developer.hashicorp.com/terraform/plugin/framework/resources/write-only-arguments) supported

## Compatibility

This provider is fully compatible with both:
- **Terraform** - The original infrastructure as code tool
- **OpenTofu** - The open-source fork of Terraform

The provider works identically with both tools and supports all the same features.

## Why

Given there already exists platform oriented, first-class providers, why do I create this? The reason is that most providers today are manually maintained, which means some latest features are likely not available in these first-class providers. For this case, `terraform-provider-restful` (or `opentofu-provider-restful`) can be used as your escape hatch.

Another common use case is that the platform you are currently working on do not have a Terraform/OpenTofu provider yet. In this case, you can use this provider to manage the resources for that platform.

## Requirement

`terraform-provider-restful` has following assumptions about the API:

- The API is expected to support the following HTTP methods:
    - `POST`/`PUT`: create the resource
    - `GET`: read the resource
    - `PUT`/`PATCH`/`POST`: update the resource
    - `DELETE`: remove the resource
- The API content type is `application/json`
- The resource should have a unique identifier (e.g. `/foos/foo1`).

Regarding the users, as this provider is essentially just a terraform/opentofu-wrapped API client, practitioners have to know the details of the API for the target platform quite well.

## Azure Resource Naming Support

The provider includes a `restful_name` data source that generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This functionality is similar to the `azurecaf_name` data source from terraform-provider-azurecaf.

### Example Usage

```hcl
data "restful_name" "example" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod", "eastus"]
  suffixes      = ["web"]
  random_length = 4
}

output "resource_group_name" {
  value = data.restful_name.example.result
  # Output: prod-eastus-rg-myapp-web-abc1
}
```

### Supported Azure Resource Types

The provider supports over 350 Azure resource types with their respective naming conventions. Here are some commonly used ones:

| Resource Type | Slug | Min Length | Max Length | Dashes | Lowercase | Example |
|---------------|------|------------|------------|---------|-----------|---------|
| azurerm_resource_group | `rg` | 1 | 90 | ✓ | ✗ | `prod-rg-myapp` |
| azurerm_storage_account | `st` | 3 | 24 | ✗ | ✓ | `prodstmyapp1234` |
| azurerm_virtual_network | `vnet` | 2 | 64 | ✓ | ✗ | `prod-vnet-myapp` |
| azurerm_subnet | `snet` | 1 | 80 | ✓ | ✗ | `prod-snet-myapp` |
| azurerm_network_security_group | `nsg` | 1 | 80 | ✓ | ✗ | `prod-nsg-myapp` |
| azurerm_public_ip | `pip` | 1 | 80 | ✓ | ✗ | `prod-pip-myapp` |
| azurerm_virtual_machine | `vm` | 1 | 15 | ✓ | ✗ | `prod-vm-myapp` |
| azurerm_kubernetes_cluster | `aks` | 1 | 63 | ✓ | ✗ | `prod-aks-myapp` |
| azurerm_key_vault | `kv` | 3 | 24 | ✓ | ✗ | `prod-kv-myapp` |
| azurerm_log_analytics_workspace | `log` | 4 | 63 | ✓ | ✗ | `prod-log-myapp` |

For the complete list of supported resource types, see [RESOURCE_TYPES.md](RESOURCE_TYPES.md).

### Features

- **Comprehensive Coverage**: Over 350 Azure resource types supported
- **Validation**: Automatic validation against Azure naming rules
- **Flexibility**: Support for prefixes, suffixes, and random suffixes
- **Compliant**: Follows Azure Cloud Adoption Framework guidelines
- **Embedded**: Resource definitions are embedded in the provider binary

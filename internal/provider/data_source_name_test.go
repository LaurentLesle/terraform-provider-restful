package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/aztfmod/opentofu-provider-restful/internal/provider"
)

// Simple provider factory for testing just the restful provider
func simpleProviderFactory() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"restful": providerserver.NewProtocol6WithError(provider.New()),
	}
}

func TestAccDataSourceName_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "myapp"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_resource_group"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceName_withPrefixesAndSuffixes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_withPrefixesAndSuffixes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "myapp"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_resource_group"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "2"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "prod"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.1", "eastus"),
					resource.TestCheckResourceAttr("data.restful_name.test", "suffixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "suffixes.0", "web"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
				),
			},
		},
	})
}

func TestAccDataSourceName_storageAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_storageAccount(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "mydata"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_storage_account"),
					resource.TestCheckResourceAttr("data.restful_name.test", "random_length", "4"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Storage account names should start with "st" slug and be lowercase
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^st[a-z0-9]+$`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_virtualMachine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_virtualMachine(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "worker"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_linux_virtual_machine"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "dev"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// VM names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^dev-vm-`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_cleanInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_cleanInput(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "My-App_With Special@Characters"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_resource_group"),
					resource.TestCheckResourceAttr("data.restful_name.test", "clean_input", "true"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Should clean special characters
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^rg-[a-zA-Z0-9-]+$`)),
				),
			},
		},
	})
}

func testAccDataSourceNameConfig_basic() string {
	return `
provider "restful" {
  base_url = "https://example.com"
}

data "restful_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
}
`
}

func testAccDataSourceNameConfig_withPrefixesAndSuffixes() string {
	return `
provider "restful" {
  base_url = "https://example.com"
}

data "restful_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod", "eastus"]
  suffixes      = ["web"]
  separator     = "-"
}
`
}

func testAccDataSourceNameConfig_storageAccount() string {
	return `
provider "restful" {
  base_url = "https://example.com"
}

data "restful_name" "test" {
  name          = "mydata"
  resource_type = "azurerm_storage_account"
  random_length = 4
}
`
}

func testAccDataSourceNameConfig_virtualMachine() string {
	return `
provider "restful" {
  base_url = "https://example.com"
}

data "restful_name" "test" {
  name          = "worker"
  resource_type = "azurerm_linux_virtual_machine"
  prefixes      = ["dev"]
}
`
}

func testAccDataSourceNameConfig_cleanInput() string {
	return `
provider "restful" {
  base_url = "https://example.com"
}

data "restful_name" "test" {
  name          = "My-App_With Special@Characters"
  resource_type = "azurerm_resource_group"
  clean_input   = true
}
`
}

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test the 5 newly added Azure resources for proper name generation and validation

func TestAccDataSourceName_containerAppJob(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_containerAppJob(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "myapp"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_container_app_job"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "dev"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Container App Job names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^dev-cajob-`)),
					// Verify lowercase constraint
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^[a-z0-9-]+$`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_grafana(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_grafana(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "monitoring"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_grafana"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "prod"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Grafana names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^prod-graf-`)),
					// Verify validation pattern (letters, numbers, hyphens, underscores - must start with letter)
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_linuxFunctionApp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_linuxFunctionApp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "processor"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_linux_function_app"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "dev"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Linux Function App names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^dev-lfa-`)),
					// Verify validation pattern (letters, numbers, hyphens - start/end with alphanumeric)
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z-]*[0-9a-zA-Z]$`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_monitorWorkspace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_monitorWorkspace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "logs"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_monitor_workspace"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "prod"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Monitor Workspace names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^prod-amws-`)),
					// Verify validation pattern (letters, numbers, periods, hyphens, underscores)
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-._]*[a-zA-Z0-9_]$`)),
				),
			},
		},
	})
}

func TestAccDataSourceName_windowsFunctionApp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: simpleProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNameConfig_windowsFunctionApp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.restful_name.test", "name", "api"),
					resource.TestCheckResourceAttr("data.restful_name.test", "resource_type", "azurerm_windows_function_app"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.restful_name.test", "prefixes.0", "staging"),
					resource.TestCheckResourceAttrSet("data.restful_name.test", "result"),
					// Windows Function App names should start with prefix, then slug
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^staging-wfa-`)),
					// Verify validation pattern (letters, numbers, hyphens - start/end with alphanumeric)
					resource.TestMatchResourceAttr("data.restful_name.test", "result", regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z-]*[0-9a-zA-Z]$`)),
				),
			},
		},
	})
}

// Configuration functions for test cases

func testAccDataSourceNameConfig_containerAppJob() string {
	return `
data "restful_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_container_app_job"
  prefixes      = ["dev"]
  clean_input   = true
}
`
}

func testAccDataSourceNameConfig_grafana() string {
	return `
data "restful_name" "test" {
  name          = "monitoring"
  resource_type = "azurerm_grafana"
  prefixes      = ["prod"]
  clean_input   = true
}
`
}

func testAccDataSourceNameConfig_linuxFunctionApp() string {
	return `
data "restful_name" "test" {
  name          = "processor"
  resource_type = "azurerm_linux_function_app"
  prefixes      = ["dev"]
  clean_input   = true
}
`
}

func testAccDataSourceNameConfig_monitorWorkspace() string {
	return `
data "restful_name" "test" {
  name          = "logs"
  resource_type = "azurerm_monitor_workspace"
  prefixes      = ["prod"]
  clean_input   = true
}
`
}

func testAccDataSourceNameConfig_windowsFunctionApp() string {
	return `
data "restful_name" "test" {
  name          = "api"
  resource_type = "azurerm_windows_function_app"
  prefixes      = ["staging"]
  clean_input   = true
}
`
}

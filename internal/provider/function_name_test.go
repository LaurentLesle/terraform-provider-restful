package provider

/*
This test file provides comprehensive test coverage for the name function that generates
standardized Azure resource names. The tests cover all functionality that was previously
tested in Terraform examples to ensure no test coverage is lost when transitioning
from Terraform testing to Go unit testing.

Test Coverage:
1. Basic functionality with required parameters (name, resource_type)
2. Singular parameter forms (prefix, suffix)
3. Plural parameter forms (prefixes, suffixes)
4. Mixed parameter usage (singular and plural together)
5. Different resource types (resource groups, storage accounts)
6. Parameter variations (separator, use_slug, passthrough, clean_input)
7. Error handling (missing required params, invalid values)
8. Edge cases (null values, empty arrays, invalid resource types)
9. Helper function extractStringArray with various input types

All test cases correspond to the working Terraform examples:
- examples/functions/name.tf
- examples/functions/simple_test.tf
- examples/functions/test_params.tf

This ensures that the Go test coverage matches the functionality proven to work
in the Terraform configuration files.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

// createDynamicValueFromJSON creates a dynamic value from JSON string for testing
func createDynamicValueFromJSON(jsonStr string) (basetypes.DynamicValue, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return basetypes.DynamicValue{}, err
	}

	value, err := convertToAttrValue(data)
	if err != nil {
		return basetypes.DynamicValue{}, err
	}

	return types.DynamicValue(value), nil
}

// convertToAttrValue converts a Go value to an attr.Value
func convertToAttrValue(data interface{}) (attr.Value, error) {
	switch v := data.(type) {
	case string:
		return types.StringValue(v), nil
	case bool:
		return types.BoolValue(v), nil
	case float64:
		bf := big.NewFloat(v)
		return types.NumberValue(bf), nil
	case nil:
		return types.StringNull(), nil
	case []interface{}:
		elements := make([]attr.Value, len(v))
		for i, elem := range v {
			val, err := convertToAttrValue(elem)
			if err != nil {
				return nil, err
			}
			elements[i] = val
		}
		return types.ListValueMust(types.StringType, elements), nil
	case map[string]interface{}:
		attrs := make(map[string]attr.Value)
		attrTypes := make(map[string]attr.Type)

		for key, val := range v {
			attrVal, err := convertToAttrValue(val)
			if err != nil {
				return nil, err
			}
			attrs[key] = attrVal
			attrTypes[key] = attrVal.Type(context.Background())
		}

		objVal, diag := types.ObjectValue(attrTypes, attrs)
		if diag.HasError() {
			return nil, fmt.Errorf("failed to create object value: %v", diag)
		}
		return objVal, nil
	default:
		return types.StringValue(""), nil
	}
}

func TestNameFunction_Metadata(t *testing.T) {
	f := nameFunction{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), function.MetadataRequest{}, resp)

	require.Equal(t, "name", resp.Name)
}

func TestNameFunction_Definition(t *testing.T) {
	f := nameFunction{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), function.DefinitionRequest{}, resp)

	require.Equal(t, "Generate standardized Azure resource names", resp.Definition.Summary)
	require.Len(t, resp.Definition.Parameters, 1)
}

func TestNameFunction_Run(t *testing.T) {
	tests := []struct {
		name        string
		inputJSON   string
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:      "basic required parameters only",
			inputJSON: `{"name": "test", "resource_type": "azurerm_resource_group"}`,
			expected:  "rg-test",
		},
		{
			name:      "storage account with suffixes array",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_storage_account", "suffixes": ["001"], "separator": "", "use_slug": true}`,
			expected:  "stmyapp001",
		},
		{
			name:      "resource group with prefix and suffix (singular)",
			inputJSON: `{"name": "myproject", "resource_type": "azurerm_resource_group", "prefix": "dev", "suffix": "001"}`,
			expected:  "dev-rg-myproject-001",
		},
		{
			name:      "test with singular prefix/suffix",
			inputJSON: `{"name": "webapp", "resource_type": "azurerm_resource_group", "prefix": "prod", "suffix": "001"}`,
			expected:  "prod-rg-webapp-001",
		},
		{
			name:      "test with plural prefixes/suffixes",
			inputJSON: `{"name": "webapp", "resource_type": "azurerm_resource_group", "prefixes": ["prod", "team1"], "suffixes": ["001", "main"]}`,
			expected:  "prod-team1-rg-webapp-001-main",
		},
		{
			name:      "test with mixed array usage",
			inputJSON: `{"name": "api", "resource_type": "azurerm_storage_account", "prefix": "dev", "suffixes": ["001", "cache"], "separator": ""}`,
			expected:  "devstapi001cache",
		},
		{
			name:      "simple resource group with custom settings",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_resource_group", "prefixes": ["prod", "team1"], "suffixes": ["main"]}`,
			expected:  "prod-team1-rg-myapp-main",
		},
		{
			name:      "test with empty separator",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_storage_account", "prefix": "dev", "suffix": "001", "separator": ""}`,
			expected:  "devstmyapp001",
		},
		{
			name:      "test with custom separator",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_resource_group", "prefix": "dev", "suffix": "001", "separator": "_"}`,
			expected:  "dev_rg_myapp_001",
		},
		{
			name:      "test with use_slug disabled",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_resource_group", "prefix": "dev", "suffix": "001", "use_slug": false}`,
			expected:  "dev-myapp-001",
		},
		{
			name:      "test with passthrough mode",
			inputJSON: `{"name": "custom-name", "resource_type": "azurerm_resource_group", "passthrough": true}`,
			expected:  "custom-name",
		},
		{
			name:      "test with clean_input disabled",
			inputJSON: `{"name": "My_App-123", "resource_type": "azurerm_resource_group", "clean_input": false}`,
			expected:  "rg-My_App-123",
		},
		// Error cases
		{
			name:        "missing name attribute",
			inputJSON:   `{"resource_type": "azurerm_resource_group"}`,
			expectError: true,
			errorMsg:    "The 'name' attribute is required in settings",
		},
		{
			name:        "missing resource_type attribute",
			inputJSON:   `{"name": "test"}`,
			expectError: true,
			errorMsg:    "The 'resource_type' attribute is required in settings",
		},
		{
			name:        "empty name attribute",
			inputJSON:   `{"name": "", "resource_type": "azurerm_resource_group"}`,
			expectError: true,
			errorMsg:    "The 'name' attribute cannot be empty",
		},
		{
			name:        "empty resource_type attribute",
			inputJSON:   `{"name": "test", "resource_type": ""}`,
			expectError: true,
			errorMsg:    "The 'resource_type' attribute cannot be empty",
		},
		{
			name:        "null name attribute",
			inputJSON:   `{"name": null, "resource_type": "azurerm_resource_group"}`,
			expectError: true,
			errorMsg:    "The 'name' attribute must be a string",
		},
		{
			name:        "null resource_type attribute",
			inputJSON:   `{"name": "test", "resource_type": null}`,
			expectError: true,
			errorMsg:    "The 'resource_type' attribute must be a string",
		},
		{
			name:        "invalid resource_type",
			inputJSON:   `{"name": "test", "resource_type": "invalid_resource_type"}`,
			expectError: true,
			errorMsg:    "Error generating name: unsupported resource type: invalid_resource_type",
		},
		{
			name:        "unauthorized attribute",
			inputJSON:   `{"name": "myapp", "resource_type": "azurerm_storage_account", "bla": "xxx"}`,
			expectError: true,
			errorMsg:    "Unauthorized attributes found: 'bla' (not supported). Supported attributes are: name, resource_type, prefix/prefixes, suffix/suffixes, separator, clean_input, use_slug, passthrough, random_length",
		},
		{
			name:        "multiple unauthorized attributes",
			inputJSON:   `{"name": "myapp", "resource_type": "azurerm_storage_account", "bla": "xxx", "another_attr": "value"}`,
			expectError: true,
			errorMsg:    "Unauthorized attributes found:",
		},
		{
			name:      "valid random_length parameter",
			inputJSON: `{"name": "myapp", "resource_type": "azurerm_storage_account", "random_length": 3, "separator": ""}`,
			expected:  "stmyappr96",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nameFunction{}

			// Parse JSON to create dynamic value
			dynamicValue, err := createDynamicValueFromJSON(tt.inputJSON)
			require.NoError(t, err, "Failed to create dynamic value from JSON")

			// Create function arguments
			args := function.NewArgumentsData([]attr.Value{dynamicValue})

			// Create request and response
			req := function.RunRequest{Arguments: args}
			resp := &function.RunResponse{}

			// Run the function
			f.Run(context.Background(), req, resp)

			if tt.expectError {
				require.Error(t, resp.Error, "Expected error but got none")
				if tt.errorMsg != "" {
					require.Contains(t, resp.Error.Error(), tt.errorMsg)
				}
			} else {
				if resp.Error != nil {
					t.Fatalf("Unexpected error: %v", resp.Error)
				}

				// Get the result value
				result := resp.Result.Value()
				stringResult, ok := result.(basetypes.StringValue)
				require.True(t, ok, "Expected string result")
				require.Equal(t, tt.expected, stringResult.ValueString())
			}
		})
	}
}

func TestNameFunction_ExtractStringArray(t *testing.T) {
	tests := []struct {
		name        string
		attrs       map[string]attr.Value
		singularKey string
		pluralKey   string
		expected    []string
	}{
		{
			name: "extract from plural list",
			attrs: map[string]attr.Value{
				"prefixes": types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("dev"),
					types.StringValue("team1"),
				}),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    []string{"dev", "team1"},
		},
		{
			name: "extract from singular string",
			attrs: map[string]attr.Value{
				"prefix": types.StringValue("dev"),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    []string{"dev"},
		},
		{
			name: "extract from tuple value",
			attrs: map[string]attr.Value{
				"prefixes": types.TupleValueMust([]attr.Type{types.StringType, types.StringType}, []attr.Value{
					types.StringValue("dev"),
					types.StringValue("team1"),
				}),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    []string{"dev", "team1"},
		},
		{
			name: "prefer plural over singular",
			attrs: map[string]attr.Value{
				"prefix": types.StringValue("singular"),
				"prefixes": types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("plural1"),
					types.StringValue("plural2"),
				}),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    []string{"plural1", "plural2"},
		},
		{
			name:        "no matching attributes",
			attrs:       map[string]attr.Value{},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    nil,
		},
		{
			name: "null list value",
			attrs: map[string]attr.Value{
				"prefixes": types.ListNull(types.StringType),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    nil,
		},
		{
			name: "empty list value",
			attrs: map[string]attr.Value{
				"prefixes": types.ListValueMust(types.StringType, []attr.Value{}),
			},
			singularKey: "prefix",
			pluralKey:   "prefixes",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractStringArray(tt.attrs, tt.singularKey, tt.pluralKey)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNameFunction_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		settings    interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "non-object settings",
			settings:    "invalid",
			expectError: true,
			errorMsg:    "Settings parameter must be an object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nameFunction{}

			// Create dynamic value from raw value
			var dynamicValue basetypes.DynamicValue
			switch v := tt.settings.(type) {
			case string:
				dynamicValue = types.DynamicValue(types.StringValue(v))
			default:
				t.Fatal("Unsupported test case type")
			}

			// Create function arguments
			args := function.NewArgumentsData([]attr.Value{dynamicValue})

			// Create request and response
			req := function.RunRequest{Arguments: args}
			resp := &function.RunResponse{}

			// Run the function
			f.Run(context.Background(), req, resp)

			if tt.expectError {
				require.Error(t, resp.Error, "Expected error but got none")
				if tt.errorMsg != "" {
					require.Contains(t, resp.Error.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, resp.Error, "Unexpected error: %v", resp.Error)
			}
		})
	}
}

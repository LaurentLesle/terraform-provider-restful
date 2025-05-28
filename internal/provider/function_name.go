package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure nameFunction satisfies the function.Function interface
var _ function.Function = &nameFunction{}

// nameFunction implements the Terraform provider function for Azure resource naming
type nameFunction struct{}

// Metadata returns metadata about the function
func (f nameFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "name"
}

// Definition defines the function's parameters and return type
func (f nameFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Generate standardized Azure resource names",
		MarkdownDescription: `Generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. 

This function implements functionality similar to azurecaf_name from the terraform-provider-azurecaf. It supports both singular (prefix, suffix) and plural (prefixes, suffixes) parameter names for flexibility.

## Supported Attributes

**Required:**
- ` + "`name`" + ` - The base name for the resource
- ` + "`resource_type`" + ` - The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account')

**Optional:**
- ` + "`prefix`" + ` / ` + "`prefixes`" + ` - Single prefix string or array of prefixes to prepend
- ` + "`suffix`" + ` / ` + "`suffixes`" + ` - Single suffix string or array of suffixes to append  
- ` + "`separator`" + ` - Character(s) used to separate name components (default: "-")
- ` + "`clean_input`" + ` - Whether to clean input strings according to resource rules (default: true)
- ` + "`use_slug`" + ` - Whether to include the resource type slug in the name (default: true)
- ` + "`passthrough`" + ` - Whether to skip processing and return the name as-is (default: false)
- ` + "`random_length`" + ` - Length of random suffix to append (default: 0, no random suffix)

## Examples

**Simple Example:**
` + "```hcl" + `
output "rg_name" {
  value = provider::restful::name({
    name          = "myproject"
    resource_type = "azurerm_resource_group"
  })
  # Result: "rg-myproject"
}
` + "```" + `

**Advanced Example:**
` + "```hcl" + `
output "storage_name" {
  value = provider::restful::name({
    name          = "data"
    resource_type = "azurerm_storage_account"
    prefixes      = ["prod", "team1"]
    suffixes      = ["001", "cache"]
    separator     = ""
    random_length = 3
  })
  # Result: "prodteam1stdatacache001r96"
}
` + "```" + `

The function validates all input attributes and will return an error for any unauthorized attributes.`,
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:                "settings",
				MarkdownDescription: "Configuration object containing all naming settings. Both 'name' and 'resource_type' attributes are required, all others are optional. You can use either 'prefix'/'suffix' (singular) or 'prefixes'/'suffixes' (plural) parameter names.",
			},
		},
		Return: function.StringReturn{},
	}
}

// Run executes the function logic
func (f nameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var settingsDynamic basetypes.DynamicValue

	// Get the settings parameter
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &settingsDynamic))
	if resp.Error != nil {
		return
	}

	// Convert dynamic value to object value
	settingsObj, ok := settingsDynamic.UnderlyingValue().(basetypes.ObjectValue)
	if !ok {
		resp.Error = function.NewFuncError("Settings parameter must be an object")
		return
	}

	// Extract attributes from the object
	attrs := settingsObj.Attributes()

	// Define authorized attributes
	authorizedAttrs := map[string]bool{
		"name":          true,
		"resource_type": true,
		"prefix":        true,
		"prefixes":      true,
		"suffix":        true,
		"suffixes":      true,
		"separator":     true,
		"clean_input":   true,
		"use_slug":      true,
		"passthrough":   true,
		"random_length": true,
	}

	// Check for unauthorized attributes
	var unauthorizedAttrs []string
	for attrName := range attrs {
		if !authorizedAttrs[attrName] {
			unauthorizedAttrs = append(unauthorizedAttrs, attrName)
		}
	}

	if len(unauthorizedAttrs) > 0 {
		var attrList []string
		for _, attr := range unauthorizedAttrs {
			attrList = append(attrList, fmt.Sprintf("'%s' (not supported)", attr))
		}
		resp.Error = function.NewFuncError(fmt.Sprintf("Unauthorized attributes found: %s. Supported attributes are: name, resource_type, prefix/prefixes, suffix/suffixes, separator, clean_input, use_slug, passthrough, random_length", strings.Join(attrList, ", ")))
		return
	}

	// Get name (required)
	nameAttr, nameExists := attrs["name"]
	if !nameExists {
		resp.Error = function.NewFuncError("The 'name' attribute is required in settings")
		return
	}
	nameVal, ok := nameAttr.(basetypes.StringValue)
	if !ok || nameVal.IsNull() || nameVal.IsUnknown() {
		resp.Error = function.NewFuncError("The 'name' attribute must be a string")
		return
	}
	name := nameVal.ValueString()
	if name == "" {
		resp.Error = function.NewFuncError("The 'name' attribute cannot be empty")
		return
	}

	// Get resource_type (required)
	resourceTypeAttr, resourceTypeExists := attrs["resource_type"]
	if !resourceTypeExists {
		resp.Error = function.NewFuncError("The 'resource_type' attribute is required in settings")
		return
	}
	resourceTypeVal, ok := resourceTypeAttr.(basetypes.StringValue)
	if !ok || resourceTypeVal.IsNull() || resourceTypeVal.IsUnknown() {
		resp.Error = function.NewFuncError("The 'resource_type' attribute must be a string")
		return
	}
	resourceType := resourceTypeVal.ValueString()
	if resourceType == "" {
		resp.Error = function.NewFuncError("The 'resource_type' attribute cannot be empty")
		return
	}

	// Get prefixes array (supports both "prefix"/"prefixes" parameter names and multiple value types)
	var prefixes []string = extractStringArray(attrs, "prefix", "prefixes")

	// Get suffixes array (supports both "suffix"/"suffixes" parameter names and multiple value types)
	var suffixes []string = extractStringArray(attrs, "suffix", "suffixes")

	// Get other optional parameters
	separator := "-" // default
	if separatorAttr, exists := attrs["separator"]; exists {
		if separatorVal, ok := separatorAttr.(basetypes.StringValue); ok && !separatorVal.IsNull() && !separatorVal.IsUnknown() {
			separator = separatorVal.ValueString()
		}
	}

	cleanInput := true // default
	if cleanInputAttr, exists := attrs["clean_input"]; exists {
		if cleanInputVal, ok := cleanInputAttr.(basetypes.BoolValue); ok && !cleanInputVal.IsNull() && !cleanInputVal.IsUnknown() {
			cleanInput = cleanInputVal.ValueBool()
		}
	}

	useSlug := true // default
	if useSlugAttr, exists := attrs["use_slug"]; exists {
		if useSlugVal, ok := useSlugAttr.(basetypes.BoolValue); ok && !useSlugVal.IsNull() && !useSlugVal.IsUnknown() {
			useSlug = useSlugVal.ValueBool()
		}
	}

	passthrough := false // default
	if passthroughAttr, exists := attrs["passthrough"]; exists {
		if passthroughVal, ok := passthroughAttr.(basetypes.BoolValue); ok && !passthroughVal.IsNull() && !passthroughVal.IsUnknown() {
			passthrough = passthroughVal.ValueBool()
		}
	}

	randomLength := 0 // default
	if randomLengthAttr, exists := attrs["random_length"]; exists {
		if randomLengthVal, ok := randomLengthAttr.(basetypes.NumberValue); ok && !randomLengthVal.IsNull() && !randomLengthVal.IsUnknown() {
			randomLengthFloat, _ := randomLengthVal.ValueBigFloat().Float64()
			randomLength = int(randomLengthFloat)
		}
	}

	// Generate the resource name
	opts := nameOptions{
		Name:         name,
		ResourceType: resourceType,
		Prefixes:     prefixes,
		Suffixes:     suffixes,
		Separator:    separator,
		CleanInput:   cleanInput,
		UseSlug:      useSlug,
		Passthrough:  passthrough,
		RandomLength: randomLength,
	}

	generatedName, err := generateName(opts)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("Error generating name: %s", err.Error()))
		return
	}

	resp.Result = function.NewResultData(types.StringValue(generatedName))
}

// extractStringArray extracts a string array from object attributes, supporting both singular and plural parameter names
func extractStringArray(attrs map[string]attr.Value, singularKey, pluralKey string) []string {
	var result []string

	// Helper function to process an attribute as array
	processArrayAttr := func(attrVal attr.Value) {
		if listVal, ok := attrVal.(basetypes.ListValue); ok && !listVal.IsNull() && !listVal.IsUnknown() {
			for _, elem := range listVal.Elements() {
				if strVal, ok := elem.(basetypes.StringValue); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					result = append(result, strVal.ValueString())
				}
			}
		} else if tupleVal, ok := attrVal.(basetypes.TupleValue); ok && !tupleVal.IsNull() && !tupleVal.IsUnknown() {
			// Handle tuple values (common in dynamic contexts when arrays are passed)
			for _, elem := range tupleVal.Elements() {
				if strVal, ok := elem.(basetypes.StringValue); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					result = append(result, strVal.ValueString())
				}
			}
		} else if strVal, ok := attrVal.(basetypes.StringValue); ok && !strVal.IsNull() && !strVal.IsUnknown() {
			// Single string value, treat as array with one element
			result = append(result, strVal.ValueString())
		}
	}

	// Check plural form first (e.g., "prefixes")
	if pluralAttr, exists := attrs[pluralKey]; exists {
		processArrayAttr(pluralAttr)
		return result
	}

	// Check singular form (e.g., "prefix")
	if singularAttr, exists := attrs[singularKey]; exists {
		processArrayAttr(singularAttr)
		return result
	}

	return result
}

// nameOptions holds the options for name generation
type nameOptions struct {
	Name         string
	ResourceType string
	Prefixes     []string
	Suffixes     []string
	RandomLength int
	RandomSeed   int64
	Separator    string
	CleanInput   bool
	Passthrough  bool
	UseSlug      bool
}

// generateName generates a resource name using the same logic as the data source
func generateName(opts nameOptions) (string, error) {
	// Get resource definition using existing function from data source
	resourceDef, err := getResourceDefinition(opts.ResourceType)
	if err != nil {
		return "", err
	}

	// Handle passthrough mode
	if opts.Passthrough {
		result := opts.Name
		if opts.CleanInput {
			result = cleanString(result, resourceDef)
		}
		if err := validateResourceName(result, resourceDef); err != nil {
			return "", err
		}
		return result, nil
	}

	// Build name components
	var components []string

	// Add prefixes
	for _, prefix := range opts.Prefixes {
		if opts.CleanInput {
			prefix = cleanString(prefix, resourceDef)
		}
		if prefix != "" {
			components = append(components, prefix)
		}
	}

	// Add resource slug if enabled
	if opts.UseSlug && resourceDef.Slug != "" {
		components = append(components, resourceDef.Slug)
	}

	// Add base name
	if opts.Name != "" {
		name := opts.Name
		if opts.CleanInput {
			name = cleanString(name, resourceDef)
		}
		components = append(components, name)
	}

	// Add suffixes
	for _, suffix := range opts.Suffixes {
		if opts.CleanInput {
			suffix = cleanString(suffix, resourceDef)
		}
		if suffix != "" {
			components = append(components, suffix)
		}
	}

	// Add random suffix if needed
	if opts.RandomLength > 0 {
		seed := opts.RandomSeed
		if seed == 0 {
			seed = 42 // Default seed for deterministic results in functions
		}
		randomSuffix := generateRandomString(opts.RandomLength, seed)
		components = append(components, randomSuffix)
	}

	// Join components
	result := ""
	if len(components) > 0 {
		result = components[0]
		for i := 1; i < len(components); i++ {
			result += opts.Separator + components[i]
		}
	}

	// Apply case transformation
	if resourceDef.LowerCase {
		result = strings.ToLower(result)
	}

	// Trim to max length
	if len(result) > resourceDef.MaxLength {
		result = result[:resourceDef.MaxLength]
	}

	// Validate the result
	if err := validateResourceName(result, resourceDef); err != nil {
		return "", err
	}

	return result, nil
}

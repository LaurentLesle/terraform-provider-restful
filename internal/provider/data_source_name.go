package provider

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceName implements the azurecaf_name data source functionality
type DataSourceName struct {
	p *Provider
}

var _ datasource.DataSource = &DataSourceName{}

// dataSourceNameData represents the data source configuration and computed values
type dataSourceNameData struct {
	Name         types.String `tfsdk:"name"`
	ResourceType types.String `tfsdk:"resource_type"`
	Prefixes     types.List   `tfsdk:"prefixes"`
	Suffixes     types.List   `tfsdk:"suffixes"`
	RandomLength types.Int64  `tfsdk:"random_length"`
	RandomSeed   types.Int64  `tfsdk:"random_seed"`
	Separator    types.String `tfsdk:"separator"`
	CleanInput   types.Bool   `tfsdk:"clean_input"`
	Passthrough  types.Bool   `tfsdk:"passthrough"`
	UseSlug      types.Bool   `tfsdk:"use_slug"`
	Result       types.String `tfsdk:"result"`
	ID           types.String `tfsdk:"id"`
}

// Metadata sets the data source type name
func (d *DataSourceName) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_name"
}

// Schema defines the structure and validation for the data source
func (d *DataSourceName) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The `restful_name` data source generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This implements functionality similar to azurecaf_name from the terraform-provider-azurecaf.",
		MarkdownDescription: "The `restful_name` data source generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This implements functionality similar to `azurecaf_name` from the terraform-provider-azurecaf.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The name to be used for generating the resource name.",
				MarkdownDescription: "The name to be used for generating the resource name.",
				Optional:            true,
			},
			"resource_type": schema.StringAttribute{
				Description:         "The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account').",
				MarkdownDescription: "The Azure resource type (e.g., `azurerm_resource_group`, `azurerm_storage_account`).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(getSupportedResourceTypes()...),
				},
			},
			"prefixes": schema.ListAttribute{
				Description:         "List of prefixes to be added to the resource name.",
				MarkdownDescription: "List of prefixes to be added to the resource name.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"suffixes": schema.ListAttribute{
				Description:         "List of suffixes to be added to the resource name.",
				MarkdownDescription: "List of suffixes to be added to the resource name.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"random_length": schema.Int64Attribute{
				Description:         "Length of the random string to be appended to the name. Defaults to 0.",
				MarkdownDescription: "Length of the random string to be appended to the name. Defaults to 0.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"random_seed": schema.Int64Attribute{
				Description:         "Seed for the random string generation. If not provided, current time will be used.",
				MarkdownDescription: "Seed for the random string generation. If not provided, current time will be used.",
				Optional:            true,
			},
			"separator": schema.StringAttribute{
				Description:         "Separator character to use between name components. Defaults to '-'.",
				MarkdownDescription: "Separator character to use between name components. Defaults to `-`.",
				Optional:            true,
			},
			"clean_input": schema.BoolAttribute{
				Description:         "Whether to clean input strings to remove non-compliant characters. Defaults to true.",
				MarkdownDescription: "Whether to clean input strings to remove non-compliant characters. Defaults to `true`.",
				Optional:            true,
			},
			"passthrough": schema.BoolAttribute{
				Description:         "Whether to enable passthrough mode (validation only). Defaults to false.",
				MarkdownDescription: "Whether to enable passthrough mode (validation only). Defaults to `false`.",
				Optional:            true,
			},
			"use_slug": schema.BoolAttribute{
				Description:         "Whether to include the resource type slug in the name. Defaults to true.",
				MarkdownDescription: "Whether to include the resource type slug in the name. Defaults to `true`.",
				Optional:            true,
			},
			"result": schema.StringAttribute{
				Description:         "The generated resource name.",
				MarkdownDescription: "The generated resource name.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "The ID of the naming convention object (same as the result value).",
				MarkdownDescription: "The ID of the naming convention object (same as the result value).",
				Computed:            true,
			},
		},
	}
}

// Configure sets up the data source with provider data
func (d *DataSourceName) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected providerData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.p = providerData.provider
}

// Read generates the resource name based on the provided configuration
func (d *DataSourceName) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceNameData

	// Read configuration from the request
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set defaults
	if data.Separator.IsNull() {
		data.Separator = types.StringValue("-")
	}
	if data.CleanInput.IsNull() {
		data.CleanInput = types.BoolValue(true)
	}
	if data.Passthrough.IsNull() {
		data.Passthrough = types.BoolValue(false)
	}
	if data.UseSlug.IsNull() {
		data.UseSlug = types.BoolValue(true)
	}
	if data.RandomLength.IsNull() {
		data.RandomLength = types.Int64Value(0)
	}

	// Generate the resource name
	result, err := d.generateResourceName(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Resource Name",
			fmt.Sprintf("Could not generate resource name: %s", err),
		)
		return
	}

	// Set computed values
	data.Result = types.StringValue(result)
	data.ID = types.StringValue(result)

	// Save the data to state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// generateResourceName implements the core naming logic
func (d *DataSourceName) generateResourceName(ctx context.Context, data dataSourceNameData) (string, error) {
	// Get resource definition
	resourceDef, err := getResourceDefinition(data.ResourceType.ValueString())
	if err != nil {
		return "", err
	}

	// Extract values
	name := data.Name.ValueString()
	separator := data.Separator.ValueString()
	cleanInput := data.CleanInput.ValueBool()
	passthrough := data.Passthrough.ValueBool()
	useSlug := data.UseSlug.ValueBool()
	randomLength := int(data.RandomLength.ValueInt64())

	// Handle prefixes
	var prefixes []string
	if !data.Prefixes.IsNull() {
		prefixElements := data.Prefixes.Elements()
		for _, elem := range prefixElements {
			if strVal, ok := elem.(types.String); ok {
				prefixes = append(prefixes, strVal.ValueString())
			}
		}
	}

	// Handle suffixes
	var suffixes []string
	if !data.Suffixes.IsNull() {
		suffixElements := data.Suffixes.Elements()
		for _, elem := range suffixElements {
			if strVal, ok := elem.(types.String); ok {
				suffixes = append(suffixes, strVal.ValueString())
			}
		}
	}

	// Generate random suffix if needed
	var randomSuffix string
	if randomLength > 0 {
		seed := time.Now().UnixNano()
		if !data.RandomSeed.IsNull() {
			seed = data.RandomSeed.ValueInt64()
		}
		randomSuffix = generateRandomString(randomLength, seed)
	}

	// Handle passthrough mode
	if passthrough {
		result := name
		if cleanInput {
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
	for _, prefix := range prefixes {
		if cleanInput {
			prefix = cleanString(prefix, resourceDef)
		}
		if prefix != "" {
			components = append(components, prefix)
		}
	}

	// Add resource slug if enabled
	if useSlug && resourceDef.Slug != "" {
		components = append(components, resourceDef.Slug)
	}

	// Add base name
	if name != "" {
		if cleanInput {
			name = cleanString(name, resourceDef)
		}
		components = append(components, name)
	}

	// Add suffixes
	for _, suffix := range suffixes {
		if cleanInput {
			suffix = cleanString(suffix, resourceDef)
		}
		if suffix != "" {
			components = append(components, suffix)
		}
	}

	// Add random suffix
	if randomSuffix != "" {
		components = append(components, randomSuffix)
	}

	// Join components
	result := strings.Join(components, separator)

	// Trim to max length
	if len(result) > resourceDef.MaxLength {
		result = result[:resourceDef.MaxLength]
	}

	// Apply lowercase if required
	if resourceDef.LowerCase {
		result = strings.ToLower(result)
	}

	// Validate the result
	if err := validateResourceName(result, resourceDef); err != nil {
		return "", err
	}

	return result, nil
}

// ResourceDefinition holds the naming rules for Azure resources
type ResourceDefinition struct {
	Name        string
	Slug        string
	MinLength   int
	MaxLength   int
	LowerCase   bool
	RegexFilter string
	Validation  string
}

// getResourceDefinition returns the naming rules for a given resource type
func getResourceDefinition(resourceType string) (*ResourceDefinition, error) {
	definitions := getResourceDefinitions()
	if def, exists := definitions[resourceType]; exists {
		return &def, nil
	}
	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

// getResourceDefinitions returns the mapping of resource types to their definitions
func getResourceDefinitions() map[string]ResourceDefinition {
	return map[string]ResourceDefinition{
		"azurerm_resource_group": {
			Name:        "Resource Group",
			Slug:        "rg",
			MinLength:   1,
			MaxLength:   90,
			LowerCase:   false,
			RegexFilter: `[^a-zA-Z0-9-._\(\)]`,
			Validation:  `^[a-zA-Z0-9-._\(\)]{0,89}[a-zA-Z0-9-_\(\)]$`,
		},
		"azurerm_storage_account": {
			Name:        "Storage Account",
			Slug:        "st",
			MinLength:   3,
			MaxLength:   24,
			LowerCase:   true,
			RegexFilter: `[^a-z0-9]`,
			Validation:  `^[a-z0-9]{3,24}$`,
		},
		"azurerm_key_vault": {
			Name:        "Key Vault",
			Slug:        "kv",
			MinLength:   3,
			MaxLength:   24,
			LowerCase:   false,
			RegexFilter: `[^a-zA-Z0-9-]`,
			Validation:  `^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$`,
		},
		"azurerm_container_registry": {
			Name:        "Container Registry",
			Slug:        "cr",
			MinLength:   5,
			MaxLength:   50,
			LowerCase:   true,
			RegexFilter: `[^a-zA-Z0-9]`,
			Validation:  `^[a-zA-Z0-9]{5,50}$`,
		},
		"azurerm_virtual_network": {
			Name:        "Virtual Network",
			Slug:        "vnet",
			MinLength:   2,
			MaxLength:   64,
			LowerCase:   false,
			RegexFilter: `[^a-zA-Z0-9-._]`,
			Validation:  `^[a-zA-Z0-9][a-zA-Z0-9-._]{0,62}[a-zA-Z0-9_]$`,
		},
		// Add more resource types as needed
	}
}

// getSupportedResourceTypes returns a list of supported resource types for validation
func getSupportedResourceTypes() []string {
	definitions := getResourceDefinitions()
	var types []string
	for resourceType := range definitions {
		types = append(types, resourceType)
	}
	return types
}

// cleanString removes characters that don't match the resource's regex filter
func cleanString(input string, resourceDef *ResourceDefinition) string {
	if resourceDef.RegexFilter == "" {
		return input
	}
	regex := regexp.MustCompile(resourceDef.RegexFilter)
	return regex.ReplaceAllString(input, "")
}

// validateResourceName validates the generated name against the resource's validation regex
func validateResourceName(name string, resourceDef *ResourceDefinition) error {
	if len(name) < resourceDef.MinLength {
		return fmt.Errorf("name '%s' is too short, minimum length is %d", name, resourceDef.MinLength)
	}
	if len(name) > resourceDef.MaxLength {
		return fmt.Errorf("name '%s' is too long, maximum length is %d", name, resourceDef.MaxLength)
	}
	if resourceDef.Validation != "" {
		regex := regexp.MustCompile(resourceDef.Validation)
		if !regex.MatchString(name) {
			return fmt.Errorf("name '%s' does not match validation pattern %s", name, resourceDef.Validation)
		}
	}
	return nil
}

// generateRandomString generates a random alphanumeric string of the specified length
func generateRandomString(length int, seed int64) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Create a new random source with the provided seed
	source := rand.NewSource(seed)
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}

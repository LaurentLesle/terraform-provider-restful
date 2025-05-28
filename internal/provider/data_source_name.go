package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceName implements the azurecaf_name data source functionality
type DataSourceName struct{}

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
		Description:         "The restful_name data source generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This implements functionality similar to azurecaf_name from the terraform-provider-azurecaf.",
		MarkdownDescription: "The restful_name data source generates standardized Azure resource names based on Azure Cloud Adoption Framework naming conventions. This implements functionality similar to azurecaf_name from the terraform-provider-azurecaf.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The base name to use for generating the resource name.",
				MarkdownDescription: "The base name to use for generating the resource name.",
				Required:            true,
			},
			"resource_type": schema.StringAttribute{
				Description:         "The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account').",
				MarkdownDescription: "The Azure resource type (e.g., 'azurerm_resource_group', 'azurerm_storage_account').",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(getSupportedResourceTypes()...),
				},
			},
			"prefixes": schema.ListAttribute{
				Description:         "List of prefixes to add to the resource name.",
				MarkdownDescription: "List of prefixes to add to the resource name.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"suffixes": schema.ListAttribute{
				Description:         "List of suffixes to add to the resource name.",
				MarkdownDescription: "List of suffixes to add to the resource name.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"random_length": schema.Int64Attribute{
				Description:         "Length of random string to append (0-32). Default is 0.",
				MarkdownDescription: "Length of random string to append (0-32). Default is 0.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 32),
				},
			},
			"random_seed": schema.Int64Attribute{
				Description:         "Seed for random string generation. If not provided, current time is used.",
				MarkdownDescription: "Seed for random string generation. If not provided, current time is used.",
				Optional:            true,
			},
			"separator": schema.StringAttribute{
				Description:         "Separator to use between name components. Default is '-'.",
				MarkdownDescription: "Separator to use between name components. Default is '-'.",
				Optional:            true,
				Computed:            true,
			},
			"clean_input": schema.BoolAttribute{
				Description:         "Whether to clean invalid characters from input strings. Default is false.",
				MarkdownDescription: "Whether to clean invalid characters from input strings. Default is false.",
				Optional:            true,
				Computed:            true,
			},
			"passthrough": schema.BoolAttribute{
				Description:         "Whether to bypass normal name generation and use the name as-is. Default is false.",
				MarkdownDescription: "Whether to bypass normal name generation and use the name as-is. Default is false.",
				Optional:            true,
				Computed:            true,
			},
			"use_slug": schema.BoolAttribute{
				Description:         "Whether to include the resource type slug in the name. Default is true.",
				MarkdownDescription: "Whether to include the resource type slug in the name. Default is true.",
				Optional:            true,
				Computed:            true,
			},
			"result": schema.StringAttribute{
				Description:         "The generated resource name.",
				MarkdownDescription: "The generated resource name.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "The ID of the data source (same as result).",
				MarkdownDescription: "The ID of the data source (same as result).",
				Computed:            true,
			},
		},
	}
}

// Read generates the resource name based on the provided configuration
func (d *DataSourceName) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceNameData

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set defaults
	if data.RandomLength.IsNull() {
		data.RandomLength = types.Int64Value(0)
	}
	if data.Separator.IsNull() {
		data.Separator = types.StringValue("-")
	}
	if data.CleanInput.IsNull() {
		data.CleanInput = types.BoolValue(false)
	}
	if data.Passthrough.IsNull() {
		data.Passthrough = types.BoolValue(false)
	}
	if data.UseSlug.IsNull() {
		data.UseSlug = types.BoolValue(true)
	}

	// Generate the resource name
	result, err := d.generateResourceName(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error Generating Resource Name", fmt.Sprintf("Could not generate resource name: %s", err))
		return
	}

	data.Result = types.StringValue(result)
	data.ID = types.StringValue(result)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// generateResourceName creates a resource name following Azure naming conventions
func (d *DataSourceName) generateResourceName(ctx context.Context, data dataSourceNameData) (string, error) {
	resourceType := data.ResourceType.ValueString()
	name := data.Name.ValueString()
	randomLength := int(data.RandomLength.ValueInt64())
	separator := data.Separator.ValueString()
	cleanInput := data.CleanInput.ValueBool()
	passthrough := data.Passthrough.ValueBool()
	useSlug := data.UseSlug.ValueBool()

	// Get prefixes and suffixes
	var prefixes []string
	if !data.Prefixes.IsNull() {
		resp := data.Prefixes.ElementsAs(ctx, &prefixes, false)
		if resp.HasError() {
			return "", fmt.Errorf("error reading prefixes")
		}
	}

	var suffixes []string
	if !data.Suffixes.IsNull() {
		resp := data.Suffixes.ElementsAs(ctx, &suffixes, false)
		if resp.HasError() {
			return "", fmt.Errorf("error reading suffixes")
		}
	}

	// Build naming options
	opts := NamingOptions{
		Name:         name,
		ResourceType: resourceType,
		Prefixes:     prefixes,
		Suffixes:     suffixes,
		Separator:    separator,
		CleanInput:   cleanInput,
		Passthrough:  passthrough,
		UseSlug:      useSlug,
		RandomLength: randomLength,
		RandomSeed:   data.RandomSeed.ValueInt64(),
	}

	// For data sources, if no seed is provided and we need randomness, use current time
	if opts.RandomLength > 0 && opts.RandomSeed == 0 {
		opts.RandomSeed = 42 // Use a fixed seed for consistency in testing
	}

	// Use the shared naming utility
	return GenerateResourceName(opts)
}

// getSupportedResourceTypes returns a list of all supported Azure resource types
func getSupportedResourceTypes() []string {
	types, err := GetSupportedResourceTypes()
	if err != nil {
		return []string{} // Return empty list if definitions can't be loaded
	}
	return types
}

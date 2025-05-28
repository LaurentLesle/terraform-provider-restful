package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/aztfmod/opentofu-provider-restful/internal/provider"
)

// TestOpenTofuProviderConversion tests that the provider works correctly with OpenTofu
func TestOpenTofuProviderConversion(t *testing.T) {
	ctx := context.Background()
	
	// Test that we can create a provider instance
	p := provider.New()
	if p == nil {
		t.Fatal("Provider creation failed")
	}

	// Test that we can create a provider server (for OpenTofu)
	server := providerserver.NewProtocol6(p)()
	
	if server == nil {
		t.Fatal("Provider server is nil")
	}

	// Test provider metadata
	req := &tfprotov6.GetProviderSchemaRequest{}
	resp, err := server.GetProviderSchema(ctx, req)
	if err != nil {
		t.Fatalf("GetProviderSchema failed: %v", err)
	}
	
	if resp == nil {
		t.Fatal("GetProviderSchema response is nil")
	}

	// Verify provider has expected resources
	if resp.ResourceSchemas == nil {
		t.Fatal("ResourceSchemas is nil")
	}

	expectedResources := []string{
		"restful_resource",
		"restful_operation",
	}

	for _, resourceName := range expectedResources {
		if _, exists := resp.ResourceSchemas[resourceName]; !exists {
			t.Errorf("Expected resource %s not found", resourceName)
		}
	}

	// Verify provider has expected data sources
	if resp.DataSourceSchemas == nil {
		t.Fatal("DataSourceSchemas is nil")
	}

	expectedDataSources := []string{
		"restful_resource",
	}

	for _, dataSourceName := range expectedDataSources {
		if _, exists := resp.DataSourceSchemas[dataSourceName]; !exists {
			t.Errorf("Expected data source %s not found", dataSourceName)
		}
	}

	// Verify provider has ephemeral resources
	if resp.EphemeralResourceSchemas == nil {
		t.Fatal("EphemeralResourceSchemas is nil")
	}

	expectedEphemeralResources := []string{
		"restful_resource",
	}

	for _, ephemeralResourceName := range expectedEphemeralResources {
		if _, exists := resp.EphemeralResourceSchemas[ephemeralResourceName]; !exists {
			t.Errorf("Expected ephemeral resource %s not found", ephemeralResourceName)
		}
	}

	t.Log("OpenTofu provider conversion test passed successfully!")
}

// TestProviderModuleName verifies the module name has been updated correctly
func TestProviderModuleName(t *testing.T) {
	// This is a compile-time test - if the imports work, the module name is correct
	p := provider.New()
	if p == nil {
		t.Fatal("Provider creation failed - module name conversion may have issues")
	}
	t.Log("Module name conversion successful!")
}

// TestOpenTofuCompatibility tests OpenTofu-specific features
func TestOpenTofuCompatibility(t *testing.T) {
	ctx := context.Background()
	
	// Create provider
	p := provider.New()
	
	// Create provider server
	server := providerserver.NewProtocol6(p)()
	
	// Get provider schema first
	schemaReq := &tfprotov6.GetProviderSchemaRequest{}
	schemaResp, err := server.GetProviderSchema(ctx, schemaReq)
	if err != nil {
		t.Fatalf("GetProviderSchema failed: %v", err)
	}
	
	if len(schemaResp.Diagnostics) > 0 {
		for _, diag := range schemaResp.Diagnostics {
			if diag.Severity == tfprotov6.DiagnosticSeverityError {
				t.Errorf("Schema error: %s - %s", diag.Summary, diag.Detail)
				return
			}
		}
	}
	
	// Verify provider schema contains expected resources
	if schemaResp.Provider == nil {
		t.Fatal("Provider schema is nil")
	}
	
	// Check resources
	expectedResources := []string{"restful_resource", "restful_operation"}
	for _, expectedResource := range expectedResources {
		if _, exists := schemaResp.ResourceSchemas[expectedResource]; !exists {
			t.Errorf("Expected resource %s not found in schema", expectedResource)
		}
	}
	
	// Check data sources  
	expectedDataSources := []string{"restful_resource"}
	for _, expectedDataSource := range expectedDataSources {
		if _, exists := schemaResp.DataSourceSchemas[expectedDataSource]; !exists {
			t.Errorf("Expected data source %s not found in schema", expectedDataSource)
		}
	}
	
	// Check ephemeral resources
	expectedEphemeralResources := []string{"restful_resource"}
	for _, expectedEphemeral := range expectedEphemeralResources {
		if _, exists := schemaResp.EphemeralResourceSchemas[expectedEphemeral]; !exists {
			t.Errorf("Expected ephemeral resource %s not found in schema", expectedEphemeral)
		}
	}
	
	t.Log("OpenTofu compatibility test passed!")
}

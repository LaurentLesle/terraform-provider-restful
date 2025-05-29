package provider

import (
	"encoding/json"
	"fmt"
	"sync"
)

// ResourceDefinition defines the naming rules and validation for a specific Azure resource type
type ResourceDefinition struct {
	Name            string `json:"name"`
	MinLength       int    `json:"min_length"`
	MaxLength       int    `json:"max_length"`
	ValidationRegex string `json:"validation_regex"`
	LowerCase       bool   `json:"lowercase"`
	Slug            string `json:"slug"`
	RegexFilter     string `json:"regex"`
	Description     string `json:"description"`
	Dashes          bool   `json:"dashes"`
	Scope           string `json:"scope"`
}

var (
	resourceDefinitionCache map[string]ResourceDefinition
	cacheOnce               sync.Once
)

// getResourceDefinition returns the resource definition for a given resource type
func getResourceDefinition(resourceType string) (*ResourceDefinition, error) {
	definitions, err := getResourceDefinitions()
	if err != nil {
		return nil, err
	}

	if def, exists := definitions[resourceType]; exists {
		return &def, nil
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

// getResourceDefinitions loads and caches all resource definitions
func getResourceDefinitions() (map[string]ResourceDefinition, error) {
	var err error
	cacheOnce.Do(func() {
		// Directly unmarshal the JSON since it now uses resource names as keys
		err = json.Unmarshal(ResourceDefinitionsJSON, &resourceDefinitionCache)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse resource definitions: %w", err)
	}

	return resourceDefinitionCache, nil
}

// GetSupportedResourceTypes returns a slice of all supported resource types
func GetSupportedResourceTypes() ([]string, error) {
	definitions, err := getResourceDefinitions()
	if err != nil {
		return []string{}, err
	}

	var resourceTypes []string
	for resourceType := range definitions {
		resourceTypes = append(resourceTypes, resourceType)
	}

	return resourceTypes, nil
}

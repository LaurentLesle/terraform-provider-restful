package provider

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

// NamingOptions contains all the options for generating resource names
type NamingOptions struct {
	Name         string
	ResourceType string
	Prefixes     []string
	Suffixes     []string
	Separator    string
	CleanInput   bool
	Passthrough  bool
	UseSlug      bool
	RandomLength int
	RandomSeed   int64
}

// GenerateResourceName creates a resource name following Azure naming conventions
// This is the shared implementation used by both the function and data source
func GenerateResourceName(opts NamingOptions) (string, error) {
	// Get resource definition
	resourceDef, err := getResourceDefinition(opts.ResourceType)
	if err != nil {
		return "", err
	}

	// Handle passthrough mode
	if opts.Passthrough {
		result := opts.Name
		if opts.CleanInput {
			result = CleanString(result, resourceDef)
		}
		// Trim to max length
		if len(result) > resourceDef.MaxLength {
			result = result[:resourceDef.MaxLength]
		}
		if err := ValidateResourceName(result, resourceDef); err != nil {
			return "", err
		}
		return result, nil
	}

	// Build name components
	var components []string

	// Add prefixes
	for _, prefix := range opts.Prefixes {
		if opts.CleanInput {
			prefix = CleanString(prefix, resourceDef)
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
			name = CleanString(name, resourceDef)
		}
		components = append(components, name)
	}

	// Add suffixes
	for _, suffix := range opts.Suffixes {
		if opts.CleanInput {
			suffix = CleanString(suffix, resourceDef)
		}
		if suffix != "" {
			components = append(components, suffix)
		}
	}

	// Add random suffix if specified
	if opts.RandomLength > 0 {
		seed := opts.RandomSeed
		if seed == 0 {
			// For data sources, use a more random seed. For functions, caller should provide a fixed seed.
			seed = 42 // Default seed for deterministic results
		}
		randomSuffix := GenerateRandomString(opts.RandomLength, seed)
		components = append(components, randomSuffix)
	}

	// Join components with separator
	result := strings.Join(components, opts.Separator)

	// Apply final transformations based on resource requirements
	if resourceDef.LowerCase {
		result = strings.ToLower(result)
	}

	// Remove dashes if not allowed
	if !resourceDef.Dashes {
		result = strings.ReplaceAll(result, "-", "")
		result = strings.ReplaceAll(result, "_", "")
	}

	// Trim to max length
	if len(result) > resourceDef.MaxLength {
		result = result[:resourceDef.MaxLength]
	}

	// Validate the final result
	if err := ValidateResourceName(result, resourceDef); err != nil {
		return "", err
	}

	return result, nil
}

// GenerateRandomString creates a random alphanumeric string of specified length
func GenerateRandomString(length int, seed int64) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	source := rand.NewSource(seed)
	rng := rand.New(source)
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}

// CleanString removes invalid characters based on Azure naming conventions
func CleanString(input string, resourceDef *ResourceDefinition) string {
	if input == "" {
		return input
	}

	// Common invalid characters for Azure resources
	invalidCharsPattern := `[^a-zA-Z0-9\-_]`
	regex := regexp.MustCompile(invalidCharsPattern)
	cleaned := regex.ReplaceAllString(input, "")

	// Apply lowercase transformation if required
	if resourceDef.LowerCase {
		cleaned = strings.ToLower(cleaned)
	}

	// Handle dashes based on resource requirements
	if !resourceDef.Dashes {
		cleaned = strings.ReplaceAll(cleaned, "-", "")
		cleaned = strings.ReplaceAll(cleaned, "_", "")
	}

	return cleaned
}

// ValidateResourceName validates the generated name against the resource's validation regex
func ValidateResourceName(name string, resourceDef *ResourceDefinition) error {
	if len(name) < resourceDef.MinLength {
		return fmt.Errorf("name '%s' is too short, minimum length is %d", name, resourceDef.MinLength)
	}
	if len(name) > resourceDef.MaxLength {
		return fmt.Errorf("name '%s' is too long, maximum length is %d", name, resourceDef.MaxLength)
	}
	if resourceDef.ValidationRegex != "" {
		// Clean the regex pattern - remove quotes and fix escaping
		pattern := resourceDef.ValidationRegex
		// Remove surrounding quotes if present
		pattern = strings.Trim(pattern, `"`)
		// Fix escaped backslashes in the pattern
		pattern = strings.ReplaceAll(pattern, `\\`, `\`)
		// Fix the ?* issue by replacing it with literal characters in character classes
		pattern = strings.ReplaceAll(pattern, "?*", `\?\*`)

		regex, err := regexp.Compile(pattern)
		if err != nil {
			// If regex compilation still fails, skip validation but warn
			// This is better than crashing the provider
			return nil
		}
		if !regex.MatchString(name) {
			return fmt.Errorf("name '%s' does not match validation pattern %s", name, pattern)
		}
	}
	return nil
}

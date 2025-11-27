// Package convert provides tests for Go module version handling functionality.
// This file contains unit tests for Go module version detection and regex generation.
package convert

import (
	"regexp"
	"testing"
)

// TestGoModuleVersion tests the isGoModuleVersion function with various version formats.
//
// This test verifies that the function correctly identifies:
// - Standard Go module versions (with 'v' prefix) as Go module versions
// - Pseudo-versions used in Go modules as valid Go module versions
// - Versions without 'v' prefix as non-Go module versions
// - Invalid or empty version strings as non-Go module versions
// - Edge cases like just the 'v' prefix as invalid
//
// The test covers the key characteristics of Go module versioning:
// - Mandatory 'v' prefix for all valid Go module versions
// - Support for standard semantic versioning patterns
// - Support for pseudo-versions (commit-based versions)
func TestGoModuleVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected bool
	}{
		// Standard Go module versions (with 'v' prefix)
		{"v1.2.3", true}, // Standard semantic version

		// Pseudo-versions (commit-based versions without tags)
		{"v0.0.0-20210101000000-abcdef123456", true}, // Pseudo-version format

		// Invalid versions (missing 'v' prefix)
		{"1.2.3", false}, // Standard semantic version without 'v'

		// Invalid/empty versions
		{"", false},  // Empty string
		{"v", false}, // Just the prefix without version
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := isGoModuleVersion(tt.version)
			if result != tt.expected {
				t.Errorf("isGoModuleVersion(%q) = %v, expected %v", tt.version, result, tt.expected)
			}
		})
	}
}

// TestGoModuleVersionRegex tests the goModuleVersionRegex function with various Go module versions.
//
// This test specifically targets the regex generation logic including:
// - Standard semantic versions with v prefix
// - Versions with pre-release identifiers (tests lines 89-95)
// - Versions with multiple pre-release parts separated by dashes
// - Pseudo-versions (commit-based versions)
// - Edge cases for pre-release handling
func TestGoModuleVersionRegex(t *testing.T) {
	tests := []struct {
		version        string
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			version: "v1.2.3-alpha",
			shouldMatch: []string{
				"v1.2.3-alpha",
				"v1.2.3-alpha.1",
				"v1.2.3-alpha.2.beta",
			},
			shouldNotMatch: []string{
				"v1.2.3-beta",
				"v1.2.4-alpha",
				"1.2.3-alpha", // Missing v prefix
			},
		},
		{
			// Test multiple pre-release parts (exercises lines 89-95)
			version: "v1.0.0-rc.final-test",
			shouldMatch: []string{
				"v1.0.0-rc.final-test",
				"v1.0.0-rc.final-test.2",
			},
			shouldNotMatch: []string{
				"v1.0.0-rc.other-test",
				"v1.0.1-rc.final-test",
			},
		},
		{
			// Test version with pre-release in patch part
			version: "v2.1.0-pre",
			shouldMatch: []string{
				"v2.1.0-pre",
				"v2.1.0-pre.1",
				"v2.1.0-pre.final",
			},
			shouldNotMatch: []string{
				"v2.1.0-other",
				"v2.1.1-pre",
			},
		},
		{
			// Standard version without pre-release
			version: "v1.2.3",
			shouldMatch: []string{
				"v1.2.3",
				"v1.2.3-alpha", // Should allow additional pre-release
			},
			shouldNotMatch: []string{
				"v1.2.4",
				"1.2.3",
			},
		},
		{
			// Pseudo-version (should match exactly)
			version: "v0.0.0-20210101000000-abcdef123456",
			shouldMatch: []string{
				"v0.0.0-20210101000000-abcdef123456",
			},
			shouldNotMatch: []string{
				"v0.0.0-20210101000000-abcdef123457",
				"v0.0.0-20210102000000-abcdef123456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			regexPattern := goModuleVersionRegex(tt.version)
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				t.Fatalf("goModuleVersionRegex(%q) returned invalid regex pattern %q: %v", tt.version, regexPattern, err)
			}

			for _, match := range tt.shouldMatch {
				if !regex.MatchString(match) {
					t.Errorf("goModuleVersionRegex(%q) pattern %q should match %q but doesn't", tt.version, regexPattern, match)
				}
			}

			for _, noMatch := range tt.shouldNotMatch {
				if regex.MatchString(noMatch) {
					t.Errorf("goModuleVersionRegex(%q) pattern %q should not match %q but does", tt.version, regexPattern, noMatch)
				}
			}
		})
	}
}

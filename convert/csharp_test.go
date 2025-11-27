// Package convert provides tests for C# NuGet version handling functionality.
// This file contains unit tests for C# version detection and regex generation.
package convert

import "testing"

// TestCSharpVersion tests the isCSharpVersion function with various version formats.
//
// This test verifies that the function correctly identifies:
// - 4-part versions (major.minor.patch.build) as C# versions
// - Versions with C#-specific pre-release identifiers (alpha, beta, rc, preview) as C# versions
// - Standard 3-part semantic versions as non-C# versions
// - Invalid or empty version strings as non-C# versions
func TestCSharpVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected bool
	}{
		// 4-part versions (major.minor.patch.build) - C# format
		{"1.2.3.4567", true},

		// C# pre-release patterns
		{"1.0.0-alpha", true},   // Alpha release
		{"1.0.0-beta001", true}, // Beta with build number
		{"1.0.0-rc.1", true},    // Release candidate
		{"1.0.0-preview", true}, // .NET preview release

		// Standard semantic versions - not C# specific
		{"1.2.3", false}, // Standard 3-part version

		// Invalid versions
		{"", false}, // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := isCSharpVersion(tt.version)
			if result != tt.expected {
				t.Errorf("isCSharpVersion(%q) = %v, expected %v", tt.version, result, tt.expected)
			}
		})
	}
}

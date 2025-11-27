// Package convert provides C# NuGet version handling functionality.
// This file contains functions specific to C# NuGet package versioning,
// which supports both 4-part versions and semantic versioning with pre-release identifiers.
package convert

import (
	"regexp"
	"strings"
)

// isCSharpVersion checks if a version follows C# NuGet versioning.
//
// C# NuGet supports two main version formats:
// 1. 4-part versions: major.minor.patch.build (e.g., 1.2.3.4567)
// 2. Semantic versions with pre-release identifiers (e.g., 1.0.0-alpha, 1.0.0-beta001)
//
// This function returns true if the version string matches either format.
//
// Examples:
//   - isCSharpVersion("1.2.3.4567") returns true (4-part version)
//   - isCSharpVersion("1.0.0-alpha") returns true (pre-release)
//   - isCSharpVersion("1.0.0-beta001") returns true (pre-release)
//   - isCSharpVersion("1.2.3") returns false (standard 3-part semantic version)
func isCSharpVersion(version string) bool {
	// C# can have 4 parts: major.minor.patch.build (like 1.2.3.4567)
	parts := strings.Split(version, ".")
	return len(parts) == 4 || containsCSharpPreRelease(version)
}

// containsCSharpPreRelease checks for C# pre-release patterns.
//
// C# NuGet packages commonly use specific pre-release identifiers:
// - alpha: Early development versions (e.g., 1.0.0-alpha, 1.0.0-alpha.1)
// - beta: Beta versions with possible build numbers (e.g., 1.0.0-beta001, 1.0.0-beta.2)
// - rc: Release candidate versions (e.g., 1.0.0-rc.1)
// - preview: Preview versions for .NET (e.g., 1.0.0-preview)
//
// Returns true if the version contains any of these C#-specific pre-release patterns.
//
// Examples:
//   - containsCSharpPreRelease("1.0.0-alpha") returns true
//   - containsCSharpPreRelease("1.0.0-beta001") returns true
//   - containsCSharpPreRelease("1.0.0-rc.1") returns true
//   - containsCSharpPreRelease("1.0.0-preview") returns true
//   - containsCSharpPreRelease("1.0.0-dev") returns false
func containsCSharpPreRelease(version string) bool {
	// C# pre-release: 1.0.0-alpha, 1.0.0-beta001, 1.0.0-rc.1
	return strings.Contains(version, "-alpha") ||
		strings.Contains(version, "-beta") ||
		strings.Contains(version, "-rc") ||
		strings.Contains(version, "-preview")
}

// csharpVersionRegex creates a regex for C# NuGet versions.
//
// This function generates a regular expression pattern that matches C# NuGet version formats.
// It handles both the main version parts and optional pre-release suffixes.
//
// The regex pattern includes:
// - Exact matching for specified version parts (major.minor.patch[.build])
// - Optional pre-release suffix matching for common C# patterns
// - Proper escaping of special regex characters in version strings
//
// For versions with pre-release suffixes, the pattern matches the exact suffix.
// For versions without pre-release suffixes, the pattern allows optional C# pre-release patterns.
//
// Parameters:
//   - version: The C# version string to create a regex for
//
// Returns:
//   - A regex pattern string that matches the specified version and compatible variations
//
// Examples:
//   - csharpVersionRegex("1.2.3.4567") generates pattern for exact 4-part version
//   - csharpVersionRegex("1.0.0-alpha") generates pattern for exact pre-release version
//   - csharpVersionRegex("1.2.3") generates pattern allowing optional C# pre-release suffixes
func csharpVersionRegex(version string) string {
	// C# versions can be: 1.2.3.4567, 1.0.0-alpha, 1.0.0-beta001, 1.0.0-rc.1
	pattern := REGEX_START

	// Handle pre-release
	mainVersion := version
	preRelease := ""
	if idx := strings.Index(version, "-"); idx != -1 {
		preRelease = version[idx:]
		mainVersion = version[:idx]
	}

	parts := strings.Split(mainVersion, ".")
	for i, part := range parts {
		if i > 0 {
			pattern += VERSION_DOT
		}
		pattern += regexp.QuoteMeta(part)
	}

	// Add pre-release pattern
	if preRelease != "" {
		pattern += regexp.QuoteMeta(preRelease)
	} else {
		// Allow common C# pre-release patterns
		pattern += `(?:-(?:alpha|beta|rc|preview)(?:\d+)?(?:\.\d+)?)?`
	}

	return pattern + REGEX_END
}

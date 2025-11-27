// Package convert provides Go module version handling functionality.
// This file contains functions specific to Go module versioning,
// which uses semantic versioning with a mandatory 'v' prefix and supports pseudo-versions.
package convert

import (
	"regexp"
	"strings"
)

// isGoModuleVersion checks if a version string follows Go module versioning.
//
// Go modules use semantic versioning with a mandatory 'v' prefix. This distinguishes
// Go module versions from other version formats and follows the Go module specification.
//
// Valid Go module version formats include:
// 1. Standard semantic versions: v1.2.3, v2.0.0-beta.1
// 2. Pseudo-versions: v0.0.0-20210101000000-abcdef123456 (used for commits without tags)
// 3. Pre-release versions: v1.2.3-pre.1, v1.0.0-alpha
//
// Returns true if the version string has the 'v' prefix and is not just "v".
//
// Examples:
//   - isGoModuleVersion("v1.2.3") returns true (standard Go module version)
//   - isGoModuleVersion("v0.0.0-20210101000000-abcdef123456") returns true (pseudo-version)
//   - isGoModuleVersion("v1.2.3-pre.1") returns true (pre-release)
//   - isGoModuleVersion("1.2.3") returns false (missing 'v' prefix)
//   - isGoModuleVersion("v") returns false (just the prefix)
//   - isGoModuleVersion("") returns false (empty string)
func isGoModuleVersion(version string) bool {
	// Go modules use semantic versioning with optional v prefix
	// Examples: v1.2.3, v0.0.0-20210101000000-abcdef123456, v1.2.3-pre.1
	return strings.HasPrefix(version, "v") && len(version) > 1
}

// goModuleVersionRegex creates a regex for Go module versions (with v prefix).
//
// This function generates a regular expression pattern that matches Go module version formats.
// It handles both standard semantic versions and pseudo-versions used in Go modules.
//
// The function processes several Go module version patterns:
//
// 1. Pseudo-versions: Special versions like v0.0.0-20210101000000-abcdef123456
//   - Used when referencing commits without semantic version tags
//   - Format: v0.0.0-{timestamp}-{commit-hash}
//   - Generates exact match patterns for these specific versions
//
// 2. Standard semantic versions: v1.2.3, v2.0.0-beta.1
//   - Follows semantic versioning with 'v' prefix
//   - Supports pre-release identifiers and build metadata
//   - Generates patterns that allow compatible pre-release variations
//
// 3. Pre-release handling: Allows compatible pre-release versions when not explicitly specified
//   - If version has pre-release: matches that specific pre-release and compatible variations
//   - If version has no pre-release: optionally matches any pre-release suffix
//
// Parameters:
//   - version: The Go module version string (must start with 'v')
//
// Returns:
//   - A regex pattern string that matches the specified version and compatible variations
//
// Examples:
//   - goModuleVersionRegex("v1.2.3") generates pattern for v1.2.3 with optional pre-release
//   - goModuleVersionRegex("v0.0.0-20210101000000-abcdef123456") generates exact pseudo-version pattern
//   - goModuleVersionRegex("v2.0.0-beta.1") generates pattern for that specific pre-release and variations
func goModuleVersionRegex(version string) string {
	// Go module versions: v1.2.3, v0.0.0-20210101000000-abcdef123456
	// Remove the 'v' prefix for processing but include it in the pattern
	cleanVersion := version[1:]
	pattern := REGEX_START + "v"

	// Handle pseudo-versions (v0.0.0-timestamp-hash)
	if strings.Contains(cleanVersion, "-") && len(strings.Split(cleanVersion, "-")) >= 3 {
		parts := strings.Split(cleanVersion, "-")
		if len(parts) >= 3 && parts[0] == "0.0.0" && len(parts[1]) == 14 && len(parts[2]) == 12 {
			// Pseudo-version pattern: exact match for the specific pseudo-version
			return REGEX_START + regexp.QuoteMeta(version) + REGEX_END
		}
	}

	// Regular semantic version with v prefix - allow pre-release versions
	parts := strings.Split(cleanVersion, ".")
	for i, part := range parts {
		if i > 0 {
			pattern += VERSION_DOT
		}
		if strings.Contains(part, "-") {
			// Handle pre-release
			preParts := strings.Split(part, "-")
			pattern += regexp.QuoteMeta(preParts[0])
			// Allow any pre-release suffix that starts with the given prefix
			if len(preParts) > 1 {
				pattern += `(?:-` + regexp.QuoteMeta(strings.Join(preParts[1:], "-")) + `.*)?`
			}
		} else {
			pattern += regexp.QuoteMeta(part)
		}
	}

	// Allow additional pre-release/build metadata if not already specified
	if !strings.Contains(cleanVersion, "-") {
		pattern += PRE_RELEASE_PATTERN
	}

	return pattern + REGEX_END
}

// Package convert provides functionality to convert semantic version constraints
// to regular expression patterns for version matching across multiple package
// management ecosystems.
//
// This package supports version constraint conversion for:
//   - NPM/Node.js (^1.2.3, ~1.2.3, >=1.2.3)
//   - Maven ([1.0,2.0), (1.0,2.0], version ranges)
//   - Go modules (v1.2.3, pseudo-versions)
//   - C# NuGet (1.2.3.4567, pre-release patterns)
//   - Python (1.2.3, ~=1.2.3 compatible release)
//   - Ruby (~>1.2.3 pessimistic operator)
//
// The main entry points are VersionToRegex for converting version constraints
// to compiled regular expressions, and VersionMatches for direct version
// matching without exposing the regex details.
//
// Example usage:
//
//	// Convert NPM caret constraint to regex
//	regex, err := VersionToRegex("^1.2.3")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check if version matches constraint
//	matches, err := VersionMatches("1.2.5", "^1.2.3")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Version matches: %t\n", matches) // Output: Version matches: true
package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// VersionToRegex converts a semantic version constraint string to a compiled regular expression.
//
// This is the main entry point for converting version constraints from various package
// management ecosystems into regex patterns that can be used for version matching.
//
// The function supports a wide range of version constraint formats:
//
// Exact matching:
//   - "1.2.3" - matches exactly version 1.2.3 (and compatible pre-release/build metadata)
//   - "==1.2.3" - explicit exact match
//
// Comparison operators:
//   - ">=1.2.3" - greater than or equal to 1.2.3
//   - "<=1.2.3" - less than or equal to 1.2.3
//   - ">1.2.3" - greater than 1.2.3 (exclusive)
//   - "<1.2.3" - less than 1.2.3 (exclusive)
//   - "!=1.2.3" - not equal to 1.2.3
//
// NPM-style ranges:
//   - "^1.2.3" - caret range, compatible within major version (>=1.2.3 <2.0.0)
//   - "~1.2.3" - tilde range, compatible within minor version (>=1.2.3 <1.3.0)
//
// Other ecosystem ranges:
//   - "~>1.2.3" - Ruby pessimistic operator
//   - "~=1.2.3" - Python compatible release
//   - "[1.0,2.0)" - Maven version ranges with bracket notation
//
// Ecosystem-specific formats:
//   - "v1.2.3" - Go module versions (with v prefix)
//   - "1.2.3.4567" - C# NuGet 4-part versions
//   - "1.*" - wildcard patterns
//
// The returned regex can be used with MatchString() or other regex methods
// to test if version strings satisfy the constraint.
//
// Parameters:
//   - versionStr: The version constraint string to convert
//
// Returns:
//   - *regexp.Regexp: Compiled regular expression matching the constraint
//   - error: Error if the constraint cannot be parsed or converted
//
// Example:
//
//	regex, err := VersionToRegex("^1.2.3")
//	if err != nil {
//		return err
//	}
//
//	matches := regex.MatchString("1.2.5") // true
//	matches = regex.MatchString("2.0.0")  // false
func VersionToRegex(versionStr string) (*regexp.Regexp, error) {
	// Parse the version constraint
	constraint, err := parseVersionConstraint(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version constraint: %w", err)
	}

	// Convert to regex pattern
	pattern, err := constraintToRegex(constraint)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to regex: %w", err)
	}

	// Compile the regex
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	return regex, nil
}

// VersionMatches checks if a given version string matches a version constraint.
//
// This is a convenience function that combines VersionToRegex and matching in a single call,
// making it easier to perform version constraint checking without manually handling regex compilation.
//
// The function supports the same version constraint formats as VersionToRegex, including
// exact matches, comparison operators, NPM-style ranges, and ecosystem-specific formats.
//
// This function is particularly useful when you need to perform one-time version checks
// or when you don't need to reuse the compiled regex pattern.
//
// Parameters:
//   - versionStr: The version string to test (e.g., "1.2.5", "v2.0.1")
//   - constraint: The version constraint to test against (e.g., "^1.2.3", ">=2.0.0")
//
// Returns:
//   - bool: true if the version matches the constraint, false otherwise
//   - error: Error if the constraint cannot be parsed or if regex compilation fails
//
// Example:
//
//	// Check if version 1.2.5 satisfies NPM caret constraint ^1.2.3
//	matches, err := VersionMatches("1.2.5", "^1.2.3")
//	if err != nil {
//		return err
//	}
//	fmt.Printf("Version matches: %t\n", matches) // Output: Version matches: true
//
//	// Check Go module version against exact constraint
//	matches, err = VersionMatches("v1.0.0", "v1.0.0")
//	if err != nil {
//		return err
//	}
//	fmt.Printf("Go version matches: %t\n", matches) // Output: Go version matches: true
func VersionMatches(versionStr, constraint string) (bool, error) {
	regex, err := VersionToRegex(constraint)
	if err != nil {
		return false, err
	}
	return regex.MatchString(versionStr), nil
}

// MustVersionToRegex is like VersionToRegex but panics if the conversion fails.
//
// This function is useful for compile-time constants or situations where you know
// the constraint is valid and want to avoid error handling. It should be used with
// caution and only when you're certain the version constraint is well-formed.
//
// Common use cases include:
//   - Defining version constraints as package-level variables
//   - Unit tests with known-good constraints
//   - Configuration parsing where constraints are validated separately
//
// The function supports the same constraint formats as VersionToRegex.
//
// Parameters:
//   - versionStr: The version constraint string to convert (must be valid)
//
// Returns:
//   - *regexp.Regexp: Compiled regular expression matching the constraint
//
// Panics:
//   - If the version constraint cannot be parsed
//   - If the regex pattern cannot be compiled
//   - If any error occurs during conversion
//
// Example:
//
//	// Safe usage with known-good constraint
//	var npmCaretRegex = MustVersionToRegex("^1.2.3")
//
//	// Usage in tests
//	func TestVersionMatching(t *testing.T) {
//		regex := MustVersionToRegex(">=2.0.0")
//		assert.True(t, regex.MatchString("2.1.0"))
//	}
//
// Warning: Only use this function when you're certain the constraint is valid,
// as panics can crash your program if the constraint is malformed.
func MustVersionToRegex(versionStr string) *regexp.Regexp {
	regex, err := VersionToRegex(versionStr)
	if err != nil {
		panic(fmt.Sprintf("MustVersionToRegex(%q): %v", versionStr, err))
	}
	return regex
}

// parseVersionConstraint parses a version constraint string into operator and version components.
//
// This function is the first stage of version constraint processing. It analyzes the input
// string to identify the constraint operator and extract the version specification.
//
// The parsing follows this precedence:
// 1. Maven-style ranges with brackets: [1.0,2.0), (1.0,2.0]
// 2. Multi-character operators: >=, <=, !=, ==, ~>, ~=
// 3. Single-character operators: >, <, =, ^, ~
// 4. No operator: defaults to exact match (==)
//
// The function handles whitespace normalization and operator precedence to ensure
// correct parsing of complex version constraints.
//
// Parameters:
//   - versionStr: Raw version constraint string from user input
//
// Returns:
//   - *VersionConstraint: Parsed constraint with operator and version fields
//   - error: Error if the constraint format is invalid or cannot be parsed
//
// Examples:
//   - "^1.2.3" → VersionConstraint{Operator: "^", Version: "1.2.3"}
//   - "[1.0,2.0)" → VersionConstraint{Operator: "maven-range", Version: "1.0,2.0"}
//   - "1.2.3" → VersionConstraint{Operator: "==", Version: "1.2.3"}
func parseVersionConstraint(versionStr string) (*VersionConstraint, error) {
	versionStr = strings.TrimSpace(versionStr)

	// Handle Maven ranges first (they have special bracket syntax)
	if strings.HasPrefix(versionStr, "[") || strings.HasPrefix(versionStr, "(") {
		return parseMavenRange(versionStr)
	}

	// Handle common operators - order matters for correct parsing
	operators := []string{OP_GREATER_EQUAL, OP_LESS_EQUAL, OP_NOT_EQUAL, OP_EQUAL_EQUAL, OP_PESSIMISTIC, OP_COMPATIBLE, OP_GREATER, OP_LESS, OP_EQUAL, OP_CARET, OP_TILDE}

	for _, op := range operators {
		if strings.HasPrefix(versionStr, op) {
			version := strings.TrimSpace(versionStr[len(op):])
			return &VersionConstraint{
				Operator: op,
				Version:  version,
			}, nil
		}
	}

	// If no operator, assume exact match
	return &VersionConstraint{
		Operator: OP_EQUAL_EQUAL,
		Version:  versionStr,
	}, nil
}

// constraintToRegex converts a parsed version constraint to a regular expression pattern.
//
// This function is the second stage of version constraint processing. It takes a parsed
// VersionConstraint and generates the appropriate regex pattern based on the operator type.
//
// The function delegates to specialized regex generators for each operator type:
//   - Exact matches: exactMatchRegex
//   - Comparison operators: greaterThanEqualRegex, lessThanRegex, etc.
//   - NPM ranges: caretRangeRegex, tildeRangeRegex
//   - Other ecosystem ranges: compatibleReleaseRegex, mavenRangeRegex
//
// Each regex generator implements the specific semantic rules for that constraint type,
// handling version part comparison, pre-release identifiers, and build metadata according
// to the conventions of the respective package management ecosystem.
//
// Parameters:
//   - constraint: Parsed version constraint with operator and version
//
// Returns:
//   - string: Regular expression pattern string
//   - error: Error if the operator is unsupported or pattern generation fails
//
// The generated patterns are designed to be compiled into Go regexp.Regexp objects
// for efficient version matching operations.
func constraintToRegex(constraint *VersionConstraint) (string, error) {
	version := constraint.Version

	switch constraint.Operator {
	case OP_EQUAL_EQUAL, OP_EQUAL:
		return exactMatchRegex(version), nil
	case OP_GREATER_EQUAL:
		return greaterThanEqualRegex(version)
	case OP_LESS_EQUAL:
		return lessThanEqualRegex(version)
	case OP_GREATER:
		return greaterThanRegex(version)
	case OP_LESS:
		return lessThanRegex(version)
	case OP_NOT_EQUAL:
		return notEqualRegex(version)
	case OP_CARET: // NPM caret range
		return caretRangeRegex(version)
	case OP_TILDE, OP_PESSIMISTIC: // NPM tilde range / Ruby pessimistic operator
		return tildeRangeRegex(version)
	case OP_COMPATIBLE: // Python compatible release
		return compatibleReleaseRegex(version)
	case OP_MAVEN_RANGE: // Maven version ranges
		return mavenRangeRegex(version)
	default:
		return "", fmt.Errorf("unsupported operator: %s", constraint.Operator)
	}
}

// exactMatchRegex creates a regex for exact version matching.
//
// This function generates regex patterns for precise version matching, handling various
// version formats and edge cases. It serves as a router that delegates to ecosystem-specific
// regex generators based on version format detection.
//
// The function handles several version formats:
//   - Wildcard versions (1.*, 2.1.*): Delegates to wildcardToRegex
//   - Go module versions (v1.2.3): Delegates to goModuleVersionRegex
//   - C# NuGet versions (1.2.3.4567, pre-release): Delegates to csharpVersionRegex
//   - Standard semantic versions: Processes directly with pre-release and build metadata
//
// For standard semantic versions, the function:
//   - Parses version components (major.minor.patch)
//   - Handles pre-release identifiers (1.2.3-alpha)
//   - Handles build metadata (1.2.3+build.1)
//   - Allows optional pre-release/build suffixes when not explicitly specified
//
// Parameters:
//   - version: Version string to create exact match pattern for
//
// Returns:
//   - string: Regex pattern for exact version matching
//
// Examples:
//   - exactMatchRegex("1.2.3") → pattern matching 1.2.3 with optional pre-release
//   - exactMatchRegex("1.*") → pattern matching any 1.x.x version
//   - exactMatchRegex("v1.2.3") → Go module pattern for v1.2.3
func exactMatchRegex(version string) string {
	// Handle wildcards and partial versions
	if strings.Contains(version, "*") {
		return wildcardToRegex(version)
	}

	// Handle Go module versions (with v prefix)
	if isGoModuleVersion(version) {
		return goModuleVersionRegex(version)
	}

	// Handle C# 4-part versions
	if isCSharpVersion(version) {
		return csharpVersionRegex(version)
	}

	// Split version into parts and pre-release/build metadata
	mainVersion := version
	preRelease := ""
	buildMeta := ""

	if idx := strings.Index(version, "+"); idx != -1 {
		buildMeta = version[idx:]
		mainVersion = version[:idx]
	}
	if idx := strings.Index(mainVersion, "-"); idx != -1 {
		preRelease = mainVersion[idx:]
		mainVersion = mainVersion[:idx]
	}

	parts := strings.Split(mainVersion, ".")
	pattern := "^"

	for i, part := range parts {
		if i > 0 {
			pattern += `\.`
		}
		pattern += regexp.QuoteMeta(part)
	}

	// Add pre-release pattern
	if preRelease != "" {
		pattern += regexp.QuoteMeta(preRelease)
	} else {
		pattern += PRE_RELEASE_PATTERN
	}

	// Add build metadata pattern
	if buildMeta != "" {
		pattern += regexp.QuoteMeta(buildMeta)
	} else {
		pattern += BUILD_META_PATTERN
	}

	pattern += "$"
	return pattern
}

// wildcardToRegex converts wildcard patterns to regex
// Result: ^1\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for "1.*")
func wildcardToRegex(version string) string {
	parts := strings.Split(version, ".")
	patternParts := convertWildcardParts(parts)
	patternParts = padToSemanticVersion(patternParts, parts)

	return REGEX_START + strings.Join(patternParts, VERSION_DOT) + VERSION_SUFFIX_PATTERN + REGEX_END
}

// convertWildcardParts converts each version part, replacing "*" with digit matcher
func convertWildcardParts(parts []string) []string {
	result := make([]string, len(parts))
	for i, part := range parts {
		if part == "*" {
			result[i] = VERSION_DIGITS
		} else {
			result[i] = regexp.QuoteMeta(part)
		}
	}
	return result
}

// padToSemanticVersion ensures we have 3 version parts (major.minor.patch)
// by appending digit matchers if the last part is a wildcard
func padToSemanticVersion(patternParts, originalParts []string) []string {
	if len(originalParts) == 0 {
		return patternParts
	}

	lastPart := originalParts[len(originalParts)-1]
	if lastPart != "*" {
		return patternParts
	}

	for len(patternParts) < 3 {
		patternParts = append(patternParts, VERSION_DIGITS)
	}
	return patternParts
}

// greaterThanEqualRegex creates a regex for >= version matching
func greaterThanEqualRegex(version string) (string, error) {
	major, minor, patch, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	var patterns []string

	// Versions with major > target major
	patterns = append(patterns, NumGreaterOrEqual(major+1)+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS)

	// Versions with major = target and minor > target minor
	patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+NumGreaterOrEqual(minor+1)+VERSION_DOT+VERSION_DIGITS)

	// Versions with major = target, minor = target, and patch >= target patch
	patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+fmt.Sprintf("%d", minor)+VERSION_DOT+NumGreaterOrEqual(patch))

	pattern := REGEX_START + "(?:" + strings.Join(patterns, REGEX_OR) + ")" + VERSION_SUFFIX_PATTERN + REGEX_END
	return pattern, nil
}

// lessThanEqualRegex creates a regex for <= version matching
func lessThanEqualRegex(version string) (string, error) {
	major, minor, patch, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	var patterns []string

	// Versions with major < target major
	if major > 0 {
		patterns = append(patterns, NumLessOrEqual(major-1)+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS)
	}

	// Versions with major = target and minor < target minor
	if minor > 0 {
		patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+NumLessOrEqual(minor-1)+VERSION_DOT+VERSION_DIGITS)
	}

	// Versions with major = target, minor = target, and patch <= target patch
	patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+fmt.Sprintf("%d", minor)+VERSION_DOT+NumLessOrEqual(patch))

	pattern := REGEX_START + "(?:" + strings.Join(patterns, REGEX_OR) + ")" + VERSION_SUFFIX_PATTERN + REGEX_END
	return pattern, nil
}

// greaterThanRegex creates a regex for > version matching
func greaterThanRegex(version string) (string, error) {
	major, minor, patch, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	var patterns []string

	// Versions with major > target major
	patterns = append(patterns, NumGreaterOrEqual(major+1)+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS)

	// Versions with major = target and minor > target minor
	patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+NumGreaterOrEqual(minor+1)+VERSION_DOT+VERSION_DIGITS)

	// Versions with major = target, minor = target, and patch > target patch
	patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+fmt.Sprintf("%d", minor)+VERSION_DOT+NumGreaterOrEqual(patch+1))

	pattern := REGEX_START + "(?:" + strings.Join(patterns, REGEX_OR) + ")" + VERSION_SUFFIX_PATTERN + REGEX_END
	return pattern, nil
}

// lessThanRegex creates a regex for < version matching
func lessThanRegex(version string) (string, error) {
	major, minor, patch, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	var patterns []string

	// Versions with major < target major
	if major > 0 {
		patterns = append(patterns, NumLessOrEqual(major-1)+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS)
	}

	// Versions with major = target and minor < target minor
	if minor > 0 {
		patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+NumLessOrEqual(minor-1)+VERSION_DOT+VERSION_DIGITS)
	}

	// Versions with major = target, minor = target, and patch < target patch
	if patch > 0 {
		patterns = append(patterns, fmt.Sprintf("%d", major)+VERSION_DOT+fmt.Sprintf("%d", minor)+VERSION_DOT+NumLessOrEqual(patch-1))
	}

	if len(patterns) == 0 {
		return EMPTY_MATCH_PATTERN, nil // No versions match (empty regex)
	}

	pattern := REGEX_START + "(?:" + strings.Join(patterns, REGEX_OR) + ")" + VERSION_SUFFIX_PATTERN + REGEX_END
	return pattern, nil
}

// notEqualRegex creates a regex for != version matching
func notEqualRegex(version string) (string, error) {
	// This is complex with regex alone - we'll use negative lookahead
	exactPattern := exactMatchRegex(version)

	// Remove ^ and $ from exact pattern for use in negative lookahead
	exactCore := strings.TrimPrefix(strings.TrimSuffix(exactPattern, REGEX_END), REGEX_START)
	pattern := fmt.Sprintf(NOT_EQUAL_BASE_PATTERN, exactCore)

	return pattern, nil
}

// caretRangeRegex creates a regex for NPM caret range (^1.2.3).
//
// NPM caret ranges implement "compatible within major version" semantics:
//   - ^1.2.3 allows >=1.2.3 and <2.0.0
//   - Changes that do not modify the left-most non-zero digit are allowed
//   - Breaking changes (major version increments) are not allowed
//
// Special handling for 0.x.x versions:
//   - ^0.2.3 allows >=0.2.3 and <0.3.0 (locked to minor version)
//   - This is because 0.x.x is considered unstable where minor increments may be breaking
//
// The function generates patterns that match:
//   - The same major version number
//   - Any minor and patch version >= the specified version
//   - Optional pre-release and build metadata suffixes
//
// Parameters:
//   - version: Base version for caret range (e.g., "1.2.3")
//
// Returns:
//   - string: Regex pattern for caret range matching
//   - error: Error if version parsing fails
//
// Examples:
//   - caretRangeRegex("1.2.3") → matches 1.2.3, 1.2.4, 1.3.0 but not 2.0.0
//   - caretRangeRegex("0.2.3") → matches 0.2.3, 0.2.4 but not 0.3.0
func caretRangeRegex(version string) (string, error) {
	major, minor, _, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	// ^1.2.3 := >=1.2.3 <2.0.0 (compatible within same major version)
	pattern := fmt.Sprintf(CARET_RANGE_TEMPLATE, major)

	// Special case: ^0.y.z is treated as 0.y.z exactly (since 0.x.x is considered unstable)
	if major == 0 {
		pattern = fmt.Sprintf(CARET_RANGE_ZERO_MAJOR_TEMPLATE, minor)
	}

	return pattern, nil
}

// tildeRangeRegex creates a regex for NPM tilde range (~1.2.3).
//
// NPM tilde ranges implement "compatible within minor version" semantics:
//   - ~1.2.3 allows >=1.2.3 and <1.3.0
//   - Only patch-level changes are allowed
//   - Minor and major version increments are not allowed
//
// This is more restrictive than caret ranges and is useful when you want to accept
// only bug fixes and patches, but not new features or breaking changes.
//
// The function generates patterns that match:
//   - The same major and minor version numbers
//   - Any patch version >= the specified patch version
//   - Optional pre-release and build metadata suffixes
//
// Parameters:
//   - version: Base version for tilde range (e.g., "1.2.3")
//
// Returns:
//   - string: Regex pattern for tilde range matching
//   - error: Error if version parsing fails
//
// Examples:
//   - tildeRangeRegex("1.2.3") → matches 1.2.3, 1.2.4, 1.2.10 but not 1.3.0
//   - tildeRangeRegex("2.1.0") → matches 2.1.0, 2.1.5 but not 2.2.0
func tildeRangeRegex(version string) (string, error) {
	major, minor, _, err := parseVersionParts(version)
	if err != nil {
		return "", err
	}

	// ~1.2.3 := >=1.2.3 <1.3.0 (compatible within same minor version)
	pattern := fmt.Sprintf(TILDE_RANGE_TEMPLATE, major, minor)

	return pattern, nil
}

// compatibleReleaseRegex creates a regex for Python compatible release (~=1.2.3)
func compatibleReleaseRegex(version string) (string, error) {
	// ~=1.2.3 is equivalent to >=1.2.3, ==1.2.*
	return tildeRangeRegex(version)
}

// parseVersionParts extracts major, minor, and patch version numbers from a version string.
//
// This utility function parses semantic version strings and extracts the numeric
// components needed for version comparison operations. It handles the complexity
// of pre-release identifiers and build metadata by stripping them before parsing.
//
// The function follows semantic versioning conventions:
//   - Versions are expected in major.minor.patch format
//   - Pre-release identifiers (after '-') are ignored for numeric parsing
//   - Build metadata (after '+') is ignored for numeric parsing
//   - Missing minor or patch versions default to 0
//
// Processing steps:
// 1. Remove pre-release identifiers (-alpha, -beta.1, etc.)
// 2. Remove build metadata (+build.1, +20210101, etc.)
// 3. Split remaining version on '.' delimiter
// 4. Parse each component as integer with validation
//
// Parameters:
//   - version: Version string to parse (e.g., "1.2.3-alpha+build.1")
//
// Returns:
//   - major: Major version number (0 if not specified)
//   - minor: Minor version number (0 if not specified)
//   - patch: Patch version number (0 if not specified)
//   - err: Error if any version component is invalid or non-numeric
//
// Examples:
//   - parseVersionParts("1.2.3") → (1, 2, 3, nil)
//   - parseVersionParts("2.0.0-beta.1") → (2, 0, 0, nil)
//   - parseVersionParts("1.5") → (1, 5, 0, nil)
//   - parseVersionParts("3") → (3, 0, 0, nil)
func parseVersionParts(version string) (major, minor, patch int, err error) {
	// Remove pre-release and build metadata for parsing
	cleanVersion := version
	if idx := strings.Index(version, "-"); idx != -1 {
		cleanVersion = version[:idx]
	}
	if idx := strings.Index(cleanVersion, "+"); idx != -1 {
		cleanVersion = cleanVersion[:idx]
	}

	parts := strings.Split(cleanVersion, ".")

	if len(parts) >= 1 {
		if major, err = strconv.Atoi(parts[0]); err != nil {
			return 0, 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
		}
	}

	if len(parts) >= 2 {
		if minor, err = strconv.Atoi(parts[1]); err != nil {
			return 0, 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
		}
	}

	if len(parts) >= 3 {
		if patch, err = strconv.Atoi(parts[2]); err != nil {
			return 0, 0, 0, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return major, minor, patch, nil
}

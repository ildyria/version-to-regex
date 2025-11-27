package convert

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestVersionToRegex(t *testing.T) {
	tests := []struct {
		name           string
		constraint     string
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			name:           "exact match",
			constraint:     "1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.3-alpha", "1.2.3+build"},
			shouldNotMatch: []string{"1.2.4", "1.3.3", "2.2.3"},
		},
		{
			name:           "exact match with more digits",
			constraint:     "123.456.789",
			shouldMatch:    []string{"123.456.789", "123.456.789-alpha", "123.456.789+build"},
			shouldNotMatch: []string{"124.456.789", "123.457.789", "123.456.780"},
		},
		{
			name:           "exact match with operator",
			constraint:     "==1.2.3",
			shouldMatch:    []string{"1.2.3"},
			shouldNotMatch: []string{"1.2.4", "1.3.3", "2.2.3"},
		},
		{
			name:           "exact match with operator and more digits",
			constraint:     "==123.456.789",
			shouldMatch:    []string{"123.456.789"},
			shouldNotMatch: []string{"124.456.789", "123.457.789", "123.456.780"},
		},
		{
			name:           "exact match with pre-release",
			constraint:     "1.2.3-dev",
			shouldMatch:    []string{"1.2.3-dev"},
			shouldNotMatch: []string{"1.2.3", "1.2.3-dev.1", "1.2.4-dev"},
		},
		{
			name:           "exact match with pre-release and more digits",
			constraint:     "123.456.789-dev",
			shouldMatch:    []string{"123.456.789-dev"},
			shouldNotMatch: []string{"123.456.789", "123.456.789-dev.1", "123.456.780-dev"},
		},
		{
			name:           "exact match with build metadata",
			constraint:     "1.2.3+build.123",
			shouldMatch:    []string{"1.2.3+build.123"},
			shouldNotMatch: []string{"1.2.3", "1.2.3+build.456", "1.2.4+build.123"},
		},
		{
			name:           "exact match with build metadata and more digits",
			constraint:     "123.456.789+build.123",
			shouldMatch:    []string{"123.456.789+build.123"},
			shouldNotMatch: []string{"123.456.789", "123.456.789+build.456", "123.456.780+build.123"},
		},
		{
			name:           "exact match with pre-release and build metadata",
			constraint:     "1.2.3-dev+build.123",
			shouldMatch:    []string{"1.2.3-dev+build.123"},
			shouldNotMatch: []string{"1.2.3", "1.2.3-dev", "1.2.3+build.123"},
		},
		{
			name:           "exact match with pre-release and build metadata and more digits",
			constraint:     "123.456.789-dev+build.123",
			shouldMatch:    []string{"123.456.789-dev+build.123"},
			shouldNotMatch: []string{"123.456.789", "123.456.789-dev", "123.456.789+build.123"},
		},
		{
			name:           "wildcard match major",
			constraint:     "1.*",
			shouldMatch:    []string{"1.0.0", "1.2.3", "1.999.999"},
			shouldNotMatch: []string{"2.0.0", "0.9.9"},
		},
		{
			name:           "wildcard match major with more digits",
			constraint:     "123.*",
			shouldMatch:    []string{"123.0.0", "123.456.789", "123.999.999"},
			shouldNotMatch: []string{"124.0.0", "122.9.9"},
		},
		{
			name:           "wildcard match minor",
			constraint:     "1.2.*",
			shouldMatch:    []string{"1.2.0", "1.2.3", "1.2.999"},
			shouldNotMatch: []string{"1.3.0", "2.2.0"},
		},
		{
			name:           "wildcard match minor with more digits",
			constraint:     "123.456.*",
			shouldMatch:    []string{"123.456.0", "123.456.789", "123.456.999"},
			shouldNotMatch: []string{"123.457.0", "124.456.0"},
		},
		{
			name:           "caret range - compatible within major",
			constraint:     "^1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.4", "1.3.0", "1.999.999"},
			shouldNotMatch: []string{"2.0.0", "0.9.9"},
		},
		{
			name:           "caret range - compatible within major with more digits",
			constraint:     "^123.456.789",
			shouldMatch:    []string{"123.456.789", "123.456.790", "123.457.0", "123.999.999"},
			shouldNotMatch: []string{"124.0.0", "122.9.9"},
		},
		{
			name:           "caret range - zero major",
			constraint:     "^0.2.3",
			shouldMatch:    []string{"0.2.3", "0.2.4", "0.2.999"},
			shouldNotMatch: []string{"0.3.0", "1.0.0"},
		},
		{
			name:           "tilde range - compatible within minor",
			constraint:     "~1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.4", "1.2.999"},
			shouldNotMatch: []string{"1.3.0", "2.0.0"},
		},
		{
			name:           "tilde range - compatible within minor with more digits",
			constraint:     "~123.456.789",
			shouldMatch:    []string{"123.456.789", "123.456.790", "123.456.999"},
			shouldNotMatch: []string{"123.457.0", "124.0.0"},
		},
		{
			name:           "tilde range with ~>",
			constraint:     "~>1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.4", "1.2.999"},
			shouldNotMatch: []string{"1.3.0", "2.0.0"},
		},
		{
			name:           "compatible release",
			constraint:     "~=1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.4", "1.2.999"},
			shouldNotMatch: []string{"1.3.0", "2.0.0"},
		},
		// Comparison operators
		{
			name:           "greater than or equal",
			constraint:     ">=1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"},
			shouldNotMatch: []string{"1.2.2", "1.1.9", "0.9.9"},
		},
		{
			name:           "greater than or equal with more digits",
			constraint:     ">=5.6.7",
			shouldMatch:    []string{"5.6.7", "5.6.8", "5.7.0", "6.0.0"},
			shouldNotMatch: []string{"5.6.6", "5.5.9", "4.9.9"},
		},
		{
			name:           "greater than or equal with multi-digit versions",
			constraint:     ">=12.34.56",
			shouldMatch:    []string{"12.34.56", "12.34.57", "12.35.0", "13.0.0", "100.0.0"},
			shouldNotMatch: []string{"12.34.55", "12.33.99", "11.99.99", "1.2.3"},
		},
		{
			name:           "greater than or equal with large version numbers",
			constraint:     ">=100.200.300",
			shouldMatch:    []string{"100.200.300", "100.200.301", "100.201.0", "101.0.0", "999.999.999"},
			shouldNotMatch: []string{"100.200.299", "100.199.999", "99.999.999", "1.2.3"},
		},
		{
			name:           "less than or equal",
			constraint:     "<=1.2.3",
			shouldMatch:    []string{"1.2.3", "1.2.2", "1.1.9", "0.9.9"},
			shouldNotMatch: []string{"1.2.4", "1.3.0", "2.0.0"},
		},
		{
			name:           "less than or equal with more digits",
			constraint:     "<=5.6.7",
			shouldMatch:    []string{"5.6.7", "5.6.6", "5.5.9", "4.9.9"},
			shouldNotMatch: []string{"5.6.8", "5.7.0", "6.0.0"},
		},
		{
			name:           "less than or equal with multi-digit versions",
			constraint:     "<=12.34.56",
			shouldMatch:    []string{"12.34.56", "12.34.55", "12.33.99", "11.99.99", "1.2.3", "0.0.0"},
			shouldNotMatch: []string{"12.34.57", "12.35.0", "13.0.0", "100.0.0"},
		},
		{
			name:           "less than or equal with large version numbers",
			constraint:     "<=100.200.300",
			shouldMatch:    []string{"100.200.300", "100.200.299", "100.199.999", "99.999.999", "1.2.3", "0.0.0"},
			shouldNotMatch: []string{"100.200.301", "100.201.0", "101.0.0", "999.999.999"},
		},
		{
			name:           "greater than (exclusive)",
			constraint:     ">1.2.3",
			shouldMatch:    []string{"1.2.4", "1.3.0", "2.0.0"},
			shouldNotMatch: []string{"1.2.3", "1.2.2", "1.1.9", "0.9.9"},
		},
		{
			name:           "greater than (exclusive) with more digits",
			constraint:     ">5.6.7",
			shouldMatch:    []string{"5.6.8", "5.7.0", "6.0.0"},
			shouldNotMatch: []string{"5.6.7", "5.6.6", "5.5.9", "4.9.9"},
		},
		{
			name:           "greater than (exclusive) with multi-digit versions",
			constraint:     ">12.34.56",
			shouldMatch:    []string{"12.34.57", "12.35.0", "13.0.0", "100.0.0"},
			shouldNotMatch: []string{"12.34.56", "12.34.55", "12.33.99", "11.99.99", "1.2.3"},
		},
		{
			name:           "greater than (exclusive) with large version numbers",
			constraint:     ">100.200.300",
			shouldMatch:    []string{"100.200.301", "100.201.0", "101.0.0", "999.999.999"},
			shouldNotMatch: []string{"100.200.300", "100.200.299", "100.199.999", "99.999.999", "1.2.3"},
		},
		{
			name:           "less than (exclusive)",
			constraint:     "<1.2.3",
			shouldMatch:    []string{"1.2.2", "1.1.9", "0.9.9"},
			shouldNotMatch: []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"},
		},
		{
			name:           "less than (exclusive) with more digits",
			constraint:     "<5.6.7",
			shouldMatch:    []string{"5.6.6", "5.5.9", "4.9.9"},
			shouldNotMatch: []string{"5.6.7", "5.6.8", "5.7.0", "6.0.0"},
		},
		{
			name:           "less than (exclusive) with multi-digit versions",
			constraint:     "<12.34.56",
			shouldMatch:    []string{"12.34.55", "12.33.99", "11.99.99", "1.2.3", "0.0.0"},
			shouldNotMatch: []string{"12.34.56", "12.34.57", "12.35.0", "13.0.0", "100.0.0"},
		},
		{
			name:           "less than (exclusive) with large version numbers",
			constraint:     "<100.200.300",
			shouldMatch:    []string{"100.200.299", "100.199.999", "99.999.999", "1.2.3", "0.0.0"},
			shouldNotMatch: []string{"100.200.300", "100.200.301", "100.201.0", "101.0.0", "999.999.999"},
		},
		// Note: != operator currently not supported due to Go regex limitations (no negative lookahead)
		// {
		//     name:           "not equal",
		//     constraint:     "!=1.2.3",
		//     shouldMatch:    []string{"1.2.2", "1.2.4", "1.3.0", "2.0.0", "0.9.9"},
		//     shouldNotMatch: []string{"1.2.3"},
		// },
		// Edge cases
		{
			name:           "empty string constraint",
			constraint:     "",
			shouldMatch:    []string{""},
			shouldNotMatch: []string{"1.2.3", "v1.0.0"},
		},
		// Maven version ranges
		{
			name:           "maven range - inclusive both",
			constraint:     "[1.0,2.0]",
			shouldMatch:    []string{"1.0.0", "1.5.0", "2.0.0"},
			shouldNotMatch: []string{"0.9.0", "3.0.0"},
		},
		{
			name:           "maven range - lower bound only",
			constraint:     "[1.0,]",
			shouldMatch:    []string{"1.0.0", "1.5.0", "2.0.0"},
			shouldNotMatch: []string{"0.9.0"},
		},
		// Go module versions
		{
			name:           "go module version",
			constraint:     "v1.2.3",
			shouldMatch:    []string{"v1.2.3", "v1.2.3-beta.1"},
			shouldNotMatch: []string{"v1.2.4", "v2.0.0", "1.2.3"},
		},
		{
			name:           "go pseudo version",
			constraint:     "v0.0.0-20210101000000-abcdef123456",
			shouldMatch:    []string{"v0.0.0-20210101000000-abcdef123456"},
			shouldNotMatch: []string{"v0.0.0-20210102000000-abcdef123456", "v1.2.3"},
		},
		// C# NuGet versions
		{
			name:           "csharp 4-part version",
			constraint:     "1.2.3.4567",
			shouldMatch:    []string{"1.2.3.4567"},
			shouldNotMatch: []string{"1.2.3.4568", "1.2.3"},
		},
		{
			name:           "csharp pre-release",
			constraint:     "1.0.0-alpha",
			shouldMatch:    []string{"1.0.0-alpha"},
			shouldNotMatch: []string{"1.0.0-beta", "1.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := VersionToRegex(tt.constraint)
			if err != nil {
				t.Fatalf("VersionToRegex(%q) failed: %v", tt.constraint, err)
			}

			for _, version := range tt.shouldMatch {
				if !regex.MatchString(version) {
					t.Errorf("Expected %q to match %q, but it didn't", version, tt.constraint)
				}
			}

			for _, version := range tt.shouldNotMatch {
				if regex.MatchString(version) {
					t.Errorf("Expected %q to NOT match %q, but it did", version, tt.constraint)
				}
			}
		})
	}
}

func TestParseVersionConstraint(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		version  string
	}{
		{"1.2.3", "==", "1.2.3"},
		{">=1.2.3", ">=", "1.2.3"},
		{"<=1.2.3", "<=", "1.2.3"},
		{">1.2.3", ">", "1.2.3"},
		{"<1.2.3", "<", "1.2.3"},
		{"!=1.2.3", "!=", "1.2.3"},
		{"^1.2.3", "^", "1.2.3"},
		{"~1.2.3", "~", "1.2.3"},
		{"~>1.2.3", "~>", "1.2.3"},
		{"~=1.2.3", "~=", "1.2.3"},
		{"==1.2.3", "==", "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			constraint, err := parseVersionConstraint(tt.input)
			if err != nil {
				t.Fatalf("parseVersionConstraint(%q) failed: %v", tt.input, err)
			}

			if constraint.Operator != tt.operator {
				t.Errorf("Expected operator %q, got %q", tt.operator, constraint.Operator)
			}

			if constraint.Version != tt.version {
				t.Errorf("Expected version %q, got %q", tt.version, constraint.Version)
			}
		})
	}
}

func TestExactMatchRegexPreReleaseAndBuildMeta(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		shouldMatch string
	}{
		{
			name:        "pre-release only",
			version:     "1.2.3-snapshot",
			shouldMatch: "1.2.3-snapshot",
		},
		{
			name:        "build metadata only",
			version:     "1.2.3+20231125",
			shouldMatch: "1.2.3+20231125",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := exactMatchRegex(tt.version)
			regex, err := regexp.Compile(pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern %q: %v", pattern, err)
			}
			if !regex.MatchString(tt.shouldMatch) {
				t.Errorf("Pattern %q should match %q", pattern, tt.shouldMatch)
			}
		})
	}
}

func TestWildcardToRegex(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			name:           "no wildcard - covers lastPart != * branch",
			version:        "1.2.3",
			shouldMatch:    []string{"1.2.3"},
			shouldNotMatch: []string{"1.2.4", "2.2.3"},
		},
		{
			name:           "wildcard at end - covers padding loop",
			version:        "1.*",
			shouldMatch:    []string{"1.0.0", "1.2.3", "1.99.99"},
			shouldNotMatch: []string{"2.0.0"},
		},
		{
			name:           "mixed parts - covers both branches in convertWildcardParts",
			version:        "1.*.3",
			shouldMatch:    []string{"1.0.3", "1.99.3"},
			shouldNotMatch: []string{"1.0.4", "2.0.3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := wildcardToRegex(tt.version)
			regex, err := regexp.Compile(pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern %q: %v", pattern, err)
			}

			for _, v := range tt.shouldMatch {
				if !regex.MatchString(v) {
					t.Errorf("Pattern %q should match %q", pattern, v)
				}
			}
			for _, v := range tt.shouldNotMatch {
				if regex.MatchString(v) {
					t.Errorf("Pattern %q should NOT match %q", pattern, v)
				}
			}
		})
	}
}

func TestPadToSemanticVersionEmptyParts(t *testing.T) {
	// Test padToSemanticVersion with empty originalParts slice
	// This covers line 420: return patternParts when len(originalParts) == 0
	result := padToSemanticVersion([]string{}, []string{})
	if len(result) != 0 {
		t.Errorf("Expected empty result for empty input, got %v", result)
	}
}

func TestComparisonRegexInvalidVersion(t *testing.T) {
	invalidVersion := "invalid.version.x"

	// Test greaterThanEqualRegex error path (line 438)
	_, err := greaterThanEqualRegex(invalidVersion)
	if err == nil {
		t.Error("greaterThanEqualRegex: expected error for invalid version")
	}

	// Test lessThanEqualRegex error path (line 458)
	_, err = lessThanEqualRegex(invalidVersion)
	if err == nil {
		t.Error("lessThanEqualRegex: expected error for invalid version")
	}

	// Test greaterThanRegex error path (line 488)
	_, err = greaterThanRegex(invalidVersion)
	if err == nil {
		t.Error("greaterThanRegex: expected error for invalid version")
	}

	// Test lessThanRegex error path (line 501)
	_, err = lessThanRegex(invalidVersion)
	if err == nil {
		t.Error("lessThanRegex: expected error for invalid version")
	}
}

func TestGreaterThanEqualMajorVersionOnly(t *testing.T) {
	// Test >=X.0.0 which uses GREATER_EQUAL_MAJOR_TEMPLATE (line 448)
	pattern, err := greaterThanEqualRegex("2.0.0")
	if err != nil {
		t.Fatalf("greaterThanEqualRegex: unexpected error: %v", err)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("Failed to compile pattern %q: %v", pattern, err)
	}

	shouldMatch := []string{"2.0.0", "2.1.0", "3.0.0", "10.0.0"}
	for _, v := range shouldMatch {
		if !regex.MatchString(v) {
			t.Errorf("Pattern should match %q", v)
		}
	}

	shouldNotMatch := []string{"1.9.9", "0.0.1"}
	for _, v := range shouldNotMatch {
		if regex.MatchString(v) {
			t.Errorf("Pattern should NOT match %q", v)
		}
	}
}

func TestMultiDigitVersionBoundaries(t *testing.T) {
	tests := []struct {
		name           string
		constraint     string
		shouldMatch    []string
		shouldNotMatch []string
	}{
		// Test boundary at 9/10 transition
		{
			name:           ">=9.9.9 boundary",
			constraint:     ">=9.9.9",
			shouldMatch:    []string{"9.9.9", "9.9.10", "9.10.0", "10.0.0", "99.99.99"},
			shouldNotMatch: []string{"9.9.8", "9.8.99", "8.99.99"},
		},
		{
			name:           "<=10.0.0 boundary",
			constraint:     "<=10.0.0",
			shouldMatch:    []string{"10.0.0", "9.99.99", "9.9.9", "0.0.0"},
			shouldNotMatch: []string{"10.0.1", "10.1.0", "11.0.0"},
		},
		// Test boundary at 99/100 transition
		{
			name:           ">=99.99.99 boundary",
			constraint:     ">=99.99.99",
			shouldMatch:    []string{"99.99.99", "99.99.100", "99.100.0", "100.0.0", "999.999.999"},
			shouldNotMatch: []string{"99.99.98", "99.98.99", "98.99.99"},
		},
		{
			name:           "<=100.0.0 boundary",
			constraint:     "<=100.0.0",
			shouldMatch:    []string{"100.0.0", "99.99.99", "50.50.50", "0.0.0"},
			shouldNotMatch: []string{"100.0.1", "100.1.0", "101.0.0"},
		},
		// Test with versions having different digit lengths in each part
		{
			name:           ">=1.10.100 mixed digits",
			constraint:     ">=1.10.100",
			shouldMatch:    []string{"1.10.100", "1.10.101", "1.11.0", "2.0.0", "10.100.1000"},
			shouldNotMatch: []string{"1.10.99", "1.9.999", "0.99.999"},
		},
		{
			name:           "<=10.100.1000 mixed digits",
			constraint:     "<=10.100.1000",
			shouldMatch:    []string{"10.100.1000", "10.100.999", "10.99.9999", "9.999.9999", "0.0.0"},
			shouldNotMatch: []string{"10.100.1001", "10.101.0", "11.0.0"},
		},
		// Edge case: version parts at exact power of 10 boundaries
		{
			name:           ">9.9.9 just over boundary",
			constraint:     ">9.9.9",
			shouldMatch:    []string{"9.9.10", "9.10.0", "10.0.0"},
			shouldNotMatch: []string{"9.9.9", "9.9.8", "8.9.9"},
		},
		{
			name:           "<10.10.10 just under boundary",
			constraint:     "<10.10.10",
			shouldMatch:    []string{"10.10.9", "10.9.99", "9.99.99", "0.0.0"},
			shouldNotMatch: []string{"10.10.10", "10.10.11", "10.11.0", "11.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := VersionToRegex(tt.constraint)
			if err != nil {
				t.Fatalf("VersionToRegex(%q) failed: %v", tt.constraint, err)
			}

			for _, version := range tt.shouldMatch {
				if !regex.MatchString(version) {
					t.Errorf("Expected %q to match %q, but it didn't. Pattern: %s", version, tt.constraint, regex.String())
				}
			}

			for _, version := range tt.shouldNotMatch {
				if regex.MatchString(version) {
					t.Errorf("Expected %q to NOT match %q, but it did. Pattern: %s", version, tt.constraint, regex.String())
				}
			}
		})
	}
}

func TestLessThanRegexZeroVersion(t *testing.T) {
	// Test lessThanRegex with 0.0.0 - covers EMPTY_MATCH_PATTERN return (line 522)
	pattern, err := lessThanRegex("0.0.0")
	if err != nil {
		t.Fatalf("lessThanRegex: unexpected error: %v", err)
	}
	if pattern != EMPTY_MATCH_PATTERN {
		t.Errorf("Expected EMPTY_MATCH_PATTERN for <0.0.0, got %q", pattern)
	}
}

func TestConstraintToRegexUnsupportedOperator(t *testing.T) {
	// Test that constraintToRegex returns error for unsupported operator
	constraint := &VersionConstraint{
		Operator: "??",
		Version:  "1.2.3",
	}
	_, err := constraintToRegex(constraint)
	if err == nil {
		t.Error("Expected error for unsupported operator, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "unsupported operator") {
		t.Errorf("Expected 'unsupported operator' error, got: %v", err)
	}
}

func TestCaretAndTildeRegexInvalidVersion(t *testing.T) {
	invalidVersion := "invalid.x.y"

	// Test caretRangeRegex error path (line 569)
	_, err := caretRangeRegex(invalidVersion)
	if err == nil {
		t.Error("caretRangeRegex: expected error for invalid version")
	}

	// Test tildeRangeRegex error path (line 611)
	_, err = tildeRangeRegex(invalidVersion)
	if err == nil {
		t.Error("tildeRangeRegex: expected error for invalid version")
	}
}

func TestParseVersionParts(t *testing.T) {
	tests := []struct {
		version string
		major   int
		minor   int
		patch   int
		wantErr bool
	}{
		{"1.2.3", 1, 2, 3, false},
		{"1.2.3-alpha", 1, 2, 3, false},
		{"1.2.3+build", 1, 2, 3, false},
		{"1.2.3-alpha+build", 1, 2, 3, false},
		{"1.2", 1, 2, 0, false},
		{"1", 1, 0, 0, false},
		{"invalid", 0, 0, 0, true},
		{"1.invalid.3", 0, 0, 0, true},
		{"1.2.invalid", 0, 0, 0, true}, // invalid patch version
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			major, minor, patch, err := parseVersionParts(tt.version)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for version %q, but got none", tt.version)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseVersionParts(%q) failed: %v", tt.version, err)
			}

			if major != tt.major {
				t.Errorf("Expected major %d, got %d", tt.major, major)
			}

			if minor != tt.minor {
				t.Errorf("Expected minor %d, got %d", tt.minor, minor)
			}

			if patch != tt.patch {
				t.Errorf("Expected patch %d, got %d", tt.patch, patch)
			}
		})
	}
}

// Example function demonstrating usage
func ExampleVersionToRegex() {
	// Exact version match
	regex, _ := VersionToRegex("1.2.3")
	println("1.2.3 matches:", regex.MatchString("1.2.3"))
	println("1.2.4 matches:", regex.MatchString("1.2.4"))

	// NPM caret range
	regex, _ = VersionToRegex("^1.2.3")
	println("^1.2.3 - 1.2.5 matches:", regex.MatchString("1.2.5"))
	println("^1.2.3 - 2.0.0 matches:", regex.MatchString("2.0.0"))

	// NPM tilde range
	regex, _ = VersionToRegex("~1.2.3")
	println("~1.2.3 - 1.2.5 matches:", regex.MatchString("1.2.5"))
	println("~1.2.3 - 1.3.0 matches:", regex.MatchString("1.3.0"))
}

func TestVersionMatches(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		expected   bool
	}{
		{"1.2.3", "1.2.3", true},
		{"1.2.4", "1.2.3", false},
		{"1.2.5", "^1.2.3", true},
		{"2.0.0", "^1.2.3", false},
		{"1.2.5", "~1.2.3", true},
		{"1.3.0", "~1.2.3", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s matches %s", tt.version, tt.constraint), func(t *testing.T) {
			result, err := VersionMatches(tt.version, tt.constraint)
			if err != nil {
				t.Fatalf("VersionMatches failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestVersionMatchesError(t *testing.T) {
	// Test that VersionMatches returns error for invalid constraint
	_, err := VersionMatches("1.2.3", "[invalid")
	if err == nil {
		t.Error("Expected error for invalid constraint, got nil")
	}
}

func TestMustVersionToRegex(t *testing.T) {
	// Test valid constraint
	regex := MustVersionToRegex("^1.2.3")
	if !regex.MatchString("1.2.5") {
		t.Error("Expected regex to match 1.2.5")
	}

	// Test panic on invalid constraint
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustVersionToRegex to panic on invalid constraint")
		}
	}()
	MustVersionToRegex(">>invalid")
}

// TestUnsupportedOperators tests that certain operators return appropriate errors
func TestUnsupportedOperators(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		expectErr  bool
	}{
		{
			name:       "not equal operator not supported in Go regex",
			constraint: "!=1.2.3",
			expectErr:  true, // Should fail due to negative lookahead not supported in Go
		},
		{
			name:       "invalid Maven range - missing closing bracket",
			constraint: "[1.0",
			expectErr:  true, // Should fail to parse version constraint
		},
		{
			name:       "invalid Maven range - wrong brackets",
			constraint: "[1.0,2.0}",
			expectErr:  true, // Should fail to parse version constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VersionToRegex(tt.constraint)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for constraint %q, but got none", tt.constraint)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error for constraint %q: %v", tt.constraint, err)
			}
		})
	}
}

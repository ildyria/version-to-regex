package convert

import (
	"regexp"
	"testing"
)

func TestParseMavenRange(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		expected string
	}{
		{"[1.0,2.0]", false, "1.0,2.0"},
		{"(1.0,2.0)", false, "1.0,2.0"},
		{"[1.0,)", false, "1.0,"},
		{"(,2.0]", false, ",2.0"},
		{"invalid", true, ""},
		{"[1.0", true, ""},
		{"[]", true, ""}, // Test length < 3 case (line 16)
		{"ab", true, ""}, // Another short string case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			constraint, err := parseMavenRange(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if constraint.Operator != "maven-range" {
				t.Errorf("Expected operator 'maven-range', got %q", constraint.Operator)
			}

			if constraint.Version != tt.expected {
				t.Errorf("Expected version %q, got %q", tt.expected, constraint.Version)
			}
		})
	}
}

// TestMavenRangeRegex tests the mavenRangeRegex function with various Maven version ranges.
//
// This test specifically targets the regex generation logic including:
// - Invalid range formats (tests line 40)
// - Same major version ranges (tests lines 55-62)
// - Different major version ranges (tests lines 59-67)
// - Lower bound only ranges (tests lines 69-78)
// - Upper bound only ranges (tests lines 79-89)
// - Edge cases and fallback patterns
func TestMavenRangeRegex(t *testing.T) {
	tests := []struct {
		name           string
		rangeStr       string
		wantErr        bool
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			name:     "invalid range format - no comma",
			rangeStr: "1.0",
			wantErr:  true,
		},
		{
			name:     "invalid range format - too many commas",
			rangeStr: "1.0,2.0,3.0",
			wantErr:  true,
		},
		{
			name:     "same major version range",
			rangeStr: "1.0,1.9",
			wantErr:  false,
			shouldMatch: []string{
				"1.0.0",
				"1.5.2",
				"1.9.9",
				"1.0.0-alpha",
				"1.2.3+build",
			},
			shouldNotMatch: []string{
				"0.9.9",
				"2.0.0",
			},
		},
		{
			name:     "different major version range",
			rangeStr: "1.0,3.0",
			wantErr:  false,
			shouldMatch: []string{
				"1.0.0",
				"2.5.2",
				"3.0.0",
				"1.0.0-alpha",
				"2.2.3+build",
			},
			shouldNotMatch: []string{
				"0.9.9",
				"4.0.0",
			},
		},
		{
			name:     "lower bound only",
			rangeStr: "2.0,",
			wantErr:  false,
			shouldMatch: []string{
				"2.0.0",
				"3.0.0",
				"10.5.2",
				"999.0.0",
				"2.0.0-alpha",
			},
			shouldNotMatch: []string{
				"1.9.9",
				"0.5.0",
			},
		},
		{
			name:     "upper bound only - positive major",
			rangeStr: ",3.0",
			wantErr:  false,
			shouldMatch: []string{
				"0.0.0",
				"1.0.0",
				"2.5.2",
				"3.0.0",
				"1.0.0-alpha",
			},
			shouldNotMatch: []string{
				"4.0.0",
				"10.0.0",
			},
		},
		{
			name:     "upper bound only - zero major",
			rangeStr: ",0.5",
			wantErr:  false,
			shouldMatch: []string{
				"0.0.0",
				"0.1.0",
				"0.5.0",
				"0.0.1-alpha",
			},
			shouldNotMatch: []string{
				"1.0.0",
				"2.0.0",
			},
		},
		{
			name:     "both bounds empty - fallback",
			rangeStr: " , ",
			wantErr:  false,
			shouldMatch: []string{
				"1.0.0",
				"2.5.2",
				"10.0.0",
				"1.0.0-alpha",
				"2.2.3+build",
			},
			shouldNotMatch: []string{
				"1.0", // Not semantic version format
				"invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexPattern, err := mavenRangeRegex(tt.rangeStr)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("mavenRangeRegex(%q) expected error but got none", tt.rangeStr)
				}
				return
			}

			if err != nil {
				t.Fatalf("mavenRangeRegex(%q) returned unexpected error: %v", tt.rangeStr, err)
			}

			// Test positive matches
			for _, match := range tt.shouldMatch {
				matched, err := regexp.MatchString(regexPattern, match)
				if err != nil {
					t.Fatalf("Invalid regex pattern %q: %v", regexPattern, err)
				}
				if !matched {
					t.Errorf("mavenRangeRegex(%q) pattern %q should match %q but doesn't",
						tt.rangeStr, regexPattern, match)
				}
			}

			// Test negative matches
			for _, noMatch := range tt.shouldNotMatch {
				matched, err := regexp.MatchString(regexPattern, noMatch)
				if err != nil {
					t.Fatalf("Invalid regex pattern %q: %v", regexPattern, err)
				}
				if matched {
					t.Errorf("mavenRangeRegex(%q) pattern %q should not match %q but does",
						tt.rangeStr, regexPattern, noMatch)
				}
			}
		})
	}
}

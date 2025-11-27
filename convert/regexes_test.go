package convert

import (
	"regexp"
	"testing"
)

func TestNumGreaterOrEqual(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			name:           "zero matches all non-negative integers",
			n:              0,
			shouldMatch:    []string{"0", "1", "9", "10", "99", "100", "999", "1000"},
			shouldNotMatch: []string{"", "a", "-1"},
		},
		{
			name:           "negative matches all non-negative integers",
			n:              -5,
			shouldMatch:    []string{"0", "1", "9", "10", "99", "100"},
			shouldNotMatch: []string{"", "a"},
		},
		{
			name:           "single digit 5",
			n:              5,
			shouldMatch:    []string{"5", "6", "7", "8", "9", "10", "99", "100", "999"},
			shouldNotMatch: []string{"0", "1", "2", "3", "4"},
		},
		{
			name:           "single digit 9",
			n:              9,
			shouldMatch:    []string{"9", "10", "11", "99", "100", "999"},
			shouldNotMatch: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"},
		},
		{
			name:           "two digits 15",
			n:              15,
			shouldMatch:    []string{"15", "16", "19", "20", "21", "99", "100", "999"},
			shouldNotMatch: []string{"0", "1", "9", "10", "11", "14"},
		},
		{
			name:           "two digits 99",
			n:              99,
			shouldMatch:    []string{"99", "100", "101", "999", "1000"},
			shouldNotMatch: []string{"0", "1", "9", "10", "50", "98"},
		},
		{
			name:           "three digits 100",
			n:              100,
			shouldMatch:    []string{"100", "101", "199", "200", "999", "1000"},
			shouldNotMatch: []string{"0", "1", "9", "10", "50", "99"},
		},
		{
			name:           "three digits 123",
			n:              123,
			shouldMatch:    []string{"123", "124", "129", "130", "199", "200", "999", "1000"},
			shouldNotMatch: []string{"0", "1", "9", "10", "99", "100", "119", "120", "121", "122"},
		},
		{
			name:           "three digits 456",
			n:              456,
			shouldMatch:    []string{"456", "457", "459", "460", "499", "500", "999", "1000"},
			shouldNotMatch: []string{"0", "1", "9", "10", "99", "100", "400", "450", "455"},
		},
		{
			name:           "four digits 1000",
			n:              1000,
			shouldMatch:    []string{"1000", "1001", "1999", "2000", "9999", "10000"},
			shouldNotMatch: []string{"0", "1", "9", "10", "99", "100", "999"},
		},
		{
			name:           "large number 12345",
			n:              12345,
			shouldMatch:    []string{"12345", "12346", "12399", "12400", "12999", "13000", "19999", "20000", "99999", "100000"},
			shouldNotMatch: []string{"0", "1", "999", "9999", "12000", "12300", "12340", "12344"},
		},
		{
			name:           "ten digits 1234567890",
			n:              1234567890,
			shouldMatch:    []string{"1234567890", "1234567891", "1234567899", "1234567900", "2000000000", "9999999999"},
			shouldNotMatch: []string{"0", "1", "999999999", "1234567889"},
		},
		{
			name:           "ten digits all nines 9999999999 (single pattern case)",
			n:              9999999999,
			shouldMatch:    []string{"9999999999"},
			shouldNotMatch: []string{"0", "1", "9999999998", "999999999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := NumGreaterOrEqual(tt.n)
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				t.Fatalf("Failed to compile regex pattern %q: %v", pattern, err)
			}

			for _, s := range tt.shouldMatch {
				if !re.MatchString(s) {
					t.Errorf("NumGreaterOrEqual(%d) pattern %q should match %q but didn't", tt.n, pattern, s)
				}
			}

			for _, s := range tt.shouldNotMatch {
				if re.MatchString(s) {
					t.Errorf("NumGreaterOrEqual(%d) pattern %q should not match %q but did", tt.n, pattern, s)
				}
			}
		})
	}
}

func TestNumLessOrEqual(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		shouldMatch    []string
		shouldNotMatch []string
		expectNoMatch  bool // For negative numbers that return a pattern matching nothing
	}{
		{
			name:          "negative matches nothing",
			n:             -1,
			expectNoMatch: true,
		},
		{
			name:           "zero matches only zero",
			n:              0,
			shouldMatch:    []string{"0"},
			shouldNotMatch: []string{"1", "2", "9", "10", "99"},
		},
		{
			name:           "single digit 5",
			n:              5,
			shouldMatch:    []string{"0", "1", "2", "3", "4", "5"},
			shouldNotMatch: []string{"6", "7", "8", "9", "10", "99"},
		},
		{
			name:           "single digit 9",
			n:              9,
			shouldMatch:    []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			shouldNotMatch: []string{"10", "11", "99", "100"},
		},
		{
			name:           "two digits 15",
			n:              15,
			shouldMatch:    []string{"0", "1", "9", "10", "11", "14", "15"},
			shouldNotMatch: []string{"16", "17", "19", "20", "99", "100"},
		},
		{
			name:           "two digits 50",
			n:              50,
			shouldMatch:    []string{"0", "1", "9", "10", "49", "50"},
			shouldNotMatch: []string{"51", "52", "99", "100"},
		},
		{
			name:           "two digits 99",
			n:              99,
			shouldMatch:    []string{"0", "1", "9", "10", "50", "98", "99"},
			shouldNotMatch: []string{"100", "101", "999"},
		},
		{
			name:           "three digits 100",
			n:              100,
			shouldMatch:    []string{"0", "1", "9", "10", "50", "99", "100"},
			shouldNotMatch: []string{"101", "102", "199", "999"},
		},
		{
			name:           "three digits 123",
			n:              123,
			shouldMatch:    []string{"0", "1", "9", "10", "99", "100", "119", "120", "121", "122", "123"},
			shouldNotMatch: []string{"124", "125", "129", "130", "199", "999"},
		},
		{
			name:           "three digits 456",
			n:              456,
			shouldMatch:    []string{"0", "1", "9", "10", "99", "100", "400", "450", "455", "456"},
			shouldNotMatch: []string{"457", "458", "459", "460", "499", "999"},
		},
		{
			name:           "four digits 1000",
			n:              1000,
			shouldMatch:    []string{"0", "1", "9", "10", "99", "100", "999", "1000"},
			shouldNotMatch: []string{"1001", "1002", "1999", "9999"},
		},
		{
			name:           "large number 12345",
			n:              12345,
			shouldMatch:    []string{"0", "1", "999", "9999", "12000", "12300", "12340", "12344", "12345"},
			shouldNotMatch: []string{"12346", "12347", "12399", "12999", "99999"},
		},
		{
			name:           "ten digits 1234567890",
			n:              1234567890,
			shouldMatch:    []string{"0", "1", "999999999", "1234567889", "1234567890"},
			shouldNotMatch: []string{"1234567891", "1234567899", "9999999999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := NumLessOrEqual(tt.n)

			// For negative numbers, the pattern returns EMPTY_MATCH_PATTERN
			// which uses Perl-style lookahead that doesn't compile in Go
			// Just verify the pattern is returned correctly
			if tt.expectNoMatch {
				if pattern != EMPTY_MATCH_PATTERN {
					t.Errorf("NumLessOrEqual(%d) expected pattern %q for no match, got %q", tt.n, EMPTY_MATCH_PATTERN, pattern)
				}
				return
			}

			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				t.Fatalf("Failed to compile regex pattern %q: %v", pattern, err)
			}

			for _, s := range tt.shouldMatch {
				if !re.MatchString(s) {
					t.Errorf("NumLessOrEqual(%d) pattern %q should match %q but didn't", tt.n, pattern, s)
				}
			}

			for _, s := range tt.shouldNotMatch {
				if re.MatchString(s) {
					t.Errorf("NumLessOrEqual(%d) pattern %q should not match %q but did", tt.n, pattern, s)
				}
			}
		})
	}
}

func TestNumGreaterOrEqualEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		wantLen int // Approximate expected pattern length to verify complexity
	}{
		{"boundary at 10", 10, 10},
		{"boundary at 100", 100, 15},
		{"boundary at 1000", 1000, 20},
		{"all nines 999", 999, 15},
		{"round number 500", 500, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := NumGreaterOrEqual(tt.n)

			// Verify pattern compiles
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				t.Fatalf("Failed to compile regex pattern %q: %v", pattern, err)
			}

			// Verify exact boundary
			if !re.MatchString(string(rune('0'+tt.n%10)) + string(rune('0'+(tt.n/10)%10))) {
				// Just verify the exact number matches
				numStr := string(rune(tt.n))
				_ = numStr // Placeholder, actual test is the boundary check below
			}

			// Verify the exact number matches
			numStr := ""
			n := tt.n
			if n == 0 {
				numStr = "0"
			} else {
				for n > 0 {
					numStr = string(rune('0'+n%10)) + numStr
					n /= 10
				}
			}
			if !re.MatchString(numStr) {
				t.Errorf("NumGreaterOrEqual(%d) should match exact number %q", tt.n, numStr)
			}

			// Verify n-1 does not match
			if tt.n > 0 {
				n = tt.n - 1
				numStr = ""
				if n == 0 {
					numStr = "0"
				} else {
					for n > 0 {
						numStr = string(rune('0'+n%10)) + numStr
						n /= 10
					}
				}
				if re.MatchString(numStr) {
					t.Errorf("NumGreaterOrEqual(%d) should not match %q", tt.n, numStr)
				}
			}
		})
	}
}

func TestNumLessOrEqualEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"boundary at 10", 10},
		{"boundary at 100", 100},
		{"boundary at 1000", 1000},
		{"all nines 999", 999},
		{"round number 500", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := NumLessOrEqual(tt.n)

			// Verify pattern compiles
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				t.Fatalf("Failed to compile regex pattern %q: %v", pattern, err)
			}

			// Helper to convert int to string
			intToStr := func(n int) string {
				if n == 0 {
					return "0"
				}
				s := ""
				for n > 0 {
					s = string(rune('0'+n%10)) + s
					n /= 10
				}
				return s
			}

			// Verify the exact number matches
			numStr := intToStr(tt.n)
			if !re.MatchString(numStr) {
				t.Errorf("NumLessOrEqual(%d) should match exact number %q", tt.n, numStr)
			}

			// Verify n+1 does not match
			numStr = intToStr(tt.n + 1)
			if re.MatchString(numStr) {
				t.Errorf("NumLessOrEqual(%d) should not match %q", tt.n, numStr)
			}
		})
	}
}

func TestRegexConstants(t *testing.T) {
	// Test that all regex constants compile successfully
	// Note: Some patterns like EMPTY_MATCH_PATTERN use Perl-style lookahead
	// which Go's regexp doesn't support, so we skip those
	constants := map[string]string{
		"VERSION_DIGITS":         VERSION_DIGITS,
		"VERSION_DOT":            VERSION_DOT,
		"PRE_RELEASE_PATTERN":    PRE_RELEASE_PATTERN,
		"BUILD_META_PATTERN":     BUILD_META_PATTERN,
		"VERSION_SUFFIX_PATTERN": VERSION_SUFFIX_PATTERN,
		"SEMANTIC_VERSION_CORE":  SEMANTIC_VERSION_CORE,
		"EXACT_VERSION_TEMPLATE": EXACT_VERSION_TEMPLATE,
	}

	for name, pattern := range constants {
		t.Run(name, func(t *testing.T) {
			_, err := regexp.Compile(pattern)
			if err != nil {
				t.Errorf("Regex constant %s with pattern %q failed to compile: %v", name, pattern, err)
			}
		})
	}
}

func TestEmptyMatchPattern(t *testing.T) {
	// EMPTY_MATCH_PATTERN uses Perl-style negative lookahead (?!.*)
	// which Go's regexp doesn't support. Just verify the constant is defined correctly.
	expected := REGEX_START + `(?!.*)` + REGEX_END
	if EMPTY_MATCH_PATTERN != expected {
		t.Errorf("EMPTY_MATCH_PATTERN = %q, want %q", EMPTY_MATCH_PATTERN, expected)
	}
}

func TestVersionDigitsPattern(t *testing.T) {
	re := regexp.MustCompile("^" + VERSION_DIGITS + "$")

	shouldMatch := []string{"0", "1", "9", "10", "99", "100", "999", "12345", "1234567890"}
	shouldNotMatch := []string{"", "a", "1.2", "-1", "1a", "a1"}

	for _, s := range shouldMatch {
		if !re.MatchString(s) {
			t.Errorf("VERSION_DIGITS should match %q", s)
		}
	}

	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("VERSION_DIGITS should not match %q", s)
		}
	}
}

func TestVersionDotPattern(t *testing.T) {
	re := regexp.MustCompile("^" + VERSION_DOT + "$")

	if !re.MatchString(".") {
		t.Error("VERSION_DOT should match '.'")
	}

	shouldNotMatch := []string{"", "a", "1", ".."}
	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("VERSION_DOT should not match %q", s)
		}
	}
}

func TestPreReleasePattern(t *testing.T) {
	re := regexp.MustCompile("^" + PRE_RELEASE_PATTERN + "$")

	shouldMatch := []string{"", "-alpha", "-beta", "-rc.1", "-alpha.1.2", "-1.0.0-alpha"}
	shouldNotMatch := []string{"-", "alpha", "-alpha beta"}

	for _, s := range shouldMatch {
		if !re.MatchString(s) {
			t.Errorf("PRE_RELEASE_PATTERN should match %q", s)
		}
	}

	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("PRE_RELEASE_PATTERN should not match %q", s)
		}
	}
}

func TestBuildMetaPattern(t *testing.T) {
	re := regexp.MustCompile("^" + BUILD_META_PATTERN + "$")

	shouldMatch := []string{"", "+build", "+build.1", "+20210101.abcdef", "+123"}
	shouldNotMatch := []string{"+", "build", "+build info"}

	for _, s := range shouldMatch {
		if !re.MatchString(s) {
			t.Errorf("BUILD_META_PATTERN should match %q", s)
		}
	}

	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("BUILD_META_PATTERN should not match %q", s)
		}
	}
}

func TestSemanticVersionCore(t *testing.T) {
	re := regexp.MustCompile("^" + SEMANTIC_VERSION_CORE + "$")

	shouldMatch := []string{"0.0.0", "1.2.3", "10.20.30", "123.456.789", "0.0.1", "1.0.0"}
	shouldNotMatch := []string{"", "1", "1.2", "1.2.3.4", "a.b.c", "1.2.3-alpha"}

	for _, s := range shouldMatch {
		if !re.MatchString(s) {
			t.Errorf("SEMANTIC_VERSION_CORE should match %q", s)
		}
	}

	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("SEMANTIC_VERSION_CORE should not match %q", s)
		}
	}
}

func TestExactVersionTemplate(t *testing.T) {
	re := regexp.MustCompile(EXACT_VERSION_TEMPLATE)

	shouldMatch := []string{
		"0.0.0",
		"1.2.3",
		"10.20.30",
		"1.2.3-alpha",
		"1.2.3-beta.1",
		"1.2.3+build",
		"1.2.3-alpha+build",
		"123.456.789",
	}
	shouldNotMatch := []string{"", "1", "1.2", "1.2.3.4", "v1.2.3"}

	for _, s := range shouldMatch {
		if !re.MatchString(s) {
			t.Errorf("EXACT_VERSION_TEMPLATE should match %q", s)
		}
	}

	for _, s := range shouldNotMatch {
		if re.MatchString(s) {
			t.Errorf("EXACT_VERSION_TEMPLATE should not match %q", s)
		}
	}
}

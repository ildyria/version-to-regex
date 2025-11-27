// Package convert provides regex constants for semantic version pattern matching.
//
// This file contains commonly used regex pattern components that are used throughout
// the version constraint conversion functionality. By centralizing these patterns as
// constants, we ensure consistency and make the code more maintainable.
package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// Basic regex pattern components for semantic versioning
const (
	// REGEX_START is the standard regex start anchor
	REGEX_START = "^"

	// REGEX_END is the standard regex end anchor
	REGEX_END = "$"

	// REGEX_OR is the alternation operator for regex patterns
	REGEX_OR = "|"

	// VERSION_DIGITS matches one or more digits in version numbers
	VERSION_DIGITS = `\d+`

	// VERSION_DOT matches a literal dot separator in version numbers
	VERSION_DOT = `\.`

	// PRE_RELEASE_PATTERN matches optional pre-release identifiers
	// Format: -alpha, -beta.1, -rc.2, etc.
	PRE_RELEASE_PATTERN = `(?:-[a-zA-Z0-9\-\.]+)?`

	// BUILD_META_PATTERN matches optional build metadata
	// Format: +build.1, +20210101.abcdef, etc.
	BUILD_META_PATTERN = `(?:\+[a-zA-Z0-9\-\.]+)?`

	// VERSION_SUFFIX_PATTERN combines pre-release and build metadata patterns
	// This is the most commonly used suffix for semantic versions
	// Result: (?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?
	VERSION_SUFFIX_PATTERN = PRE_RELEASE_PATTERN + BUILD_META_PATTERN

	// SEMANTIC_VERSION_CORE matches the core major.minor.patch pattern
	// Result: \d+\.\d+\.\d+
	SEMANTIC_VERSION_CORE = VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS
)

// Template patterns for common version constraint formats
const (
	// EXACT_VERSION_TEMPLATE matches an exact semantic version with optional suffixes
	// Result: ^\d+\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$
	EXACT_VERSION_TEMPLATE = REGEX_START + SEMANTIC_VERSION_CORE + VERSION_SUFFIX_PATTERN + REGEX_END

	// CARET_RANGE_TEMPLATE matches NPM caret range versions (^1.2.3)
	// Placeholder: %d for major version number
	// Result: ^1\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (when %d is 1)
	CARET_RANGE_TEMPLATE = REGEX_START + "%d" + VERSION_DOT + VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS + VERSION_SUFFIX_PATTERN + REGEX_END

	// CARET_RANGE_ZERO_MAJOR_TEMPLATE matches NPM caret range for 0.x.x versions
	// Placeholder: %d for minor version number
	// Result: ^0\.1\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (when %d is 1)
	CARET_RANGE_ZERO_MAJOR_TEMPLATE = REGEX_START + "0" + VERSION_DOT + "%d" + VERSION_DOT + VERSION_DIGITS + VERSION_SUFFIX_PATTERN + REGEX_END

	// TILDE_RANGE_TEMPLATE matches NPM tilde range versions (~1.2.3)
	// Placeholders: %d for major, %d for minor version numbers
	// Result: ^1\.2\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (when %d is 1, %d is 2)
	TILDE_RANGE_TEMPLATE = REGEX_START + "%d" + VERSION_DOT + "%d" + VERSION_DOT + VERSION_DIGITS + VERSION_SUFFIX_PATTERN + REGEX_END
)

// Complex comparison patterns for >= and <= operators
const (
	// GREATER_EQUAL_MAJOR_TEMPLATE matches versions with major >= specified value
	// Placeholder: %d for minimum major version
	// Result: ^(?:[1-9]|\d{2,})\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (when %d is 1)
	GREATER_EQUAL_MAJOR_TEMPLATE = REGEX_START + `(?:[%d-9]|\d{2,})` + VERSION_DOT + VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS + VERSION_SUFFIX_PATTERN + REGEX_END

	// COMPARISON_COMPLEX_TEMPLATE matches complex >= patterns with major/minor/patch alternatives
	// Placeholders: major+1, major, minor+1, major, minor, patch
	// Result: ^(?:(?:[2-9]|\d{2,})\.\d+\.\d+|1\.(?:[3-9]|\d{2,})\.\d+|1\.2\.(?:[4-9]|\d{2,}))(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for >=1.2.3)
	COMPARISON_COMPLEX_TEMPLATE = REGEX_START + `(?:(?:[%d-9]|\d{2,})` + VERSION_DOT + VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS + `|%d` + VERSION_DOT + `(?:[%d-9]|\d{2,})` + VERSION_DOT + VERSION_DIGITS + `|%d` + VERSION_DOT + `%d` + VERSION_DOT + `(?:[%d-9]|\d{2,}))` + VERSION_SUFFIX_PATTERN + REGEX_END

	// LESS_EQUAL_SUFFIX_PATTERN provides the suffix for <= comparison patterns
	// Result: (?:-[a-zA-Z0-9\\-\\.]+)?(?:\\+[a-zA-Z0-9\\-\\.]+)?$
	LESS_EQUAL_SUFFIX_PATTERN = `(?:-[a-zA-Z0-9\\-\\.]+)?(?:\\+[a-zA-Z0-9\\-\\.]+)?` + REGEX_END
)

// Special patterns for edge cases and advanced matching
const (
	// NOT_EQUAL_BASE_PATTERN provides base structure for != operations
	// Result: ^(?!1\.2\.3$)\d+(?:\.\d+)?(?:\.\d+)?(?:-[a-zA-Z0-9\\-\\.]+)?(?:\\+[a-zA-Z0-9\\-\\.]+)?$ (when %s is 1\.2\.3)
	NOT_EQUAL_BASE_PATTERN = REGEX_START + `(?!%s` + REGEX_END + `)` + `\d+(?:\.\d+)?(?:\.\d+)?` + `(?:-[a-zA-Z0-9\\-\\.]+)?(?:\\+[a-zA-Z0-9\\-\\.]+)?` + REGEX_END

	// WILDCARD_VERSION_DIGITS matches wildcard positions in version strings
	WILDCARD_VERSION_DIGITS = VERSION_DIGITS

	// EMPTY_MATCH_PATTERN matches no versions (used for impossible constraints)
	EMPTY_MATCH_PATTERN = REGEX_START + `(?!.*)` + REGEX_END
)

// NumGreaterOrEqual generates a regex pattern that matches integers >= n.
//
// This function creates regex patterns for matching version number components
// that are greater than or equal to a given value. It supports multi-digit
// numbers up to 10 digits (0 to 9999999999), making it suitable for semantic
// version comparisons where version parts can be arbitrarily large.
//
// Algorithm:
// The function builds multiple alternation patterns to cover all cases:
//  1. Numbers with more digits than n (always greater, e.g., 100+ when n=15)
//  2. Numbers with the same digit count but higher leading digits
//  3. Numbers matching the exact prefix with higher subsequent digits
//
// Pattern Construction Strategy:
//   - For n=0: matches any non-negative integer (\d+)
//   - For single digits (n=5): matches [5-9] or multi-digit numbers
//   - For multi-digit (n=15): matches 15-19 (1[5-9]), 20-99 ([2-9]\d), or 100+ (\d{3,})
//
// Parameters:
//   - n: The minimum value to match (inclusive). Must be >= 0.
//
// Returns:
//   - A regex pattern string that matches any integer >= n.
//     The pattern uses non-capturing groups (?:...) for alternation.
//
// Examples:
//
//	NumGreaterOrEqual(0)   → `\d+`
//	    Matches: 0, 1, 10, 100, 9999999999
//
//	NumGreaterOrEqual(5)   → `(?:[5-9]|\d{2,})`
//	    Matches: 5, 6, 7, 8, 9, 10, 100, 1000...
//	    Does not match: 0, 1, 2, 3, 4
//
//	NumGreaterOrEqual(15)  → `(?:1[5-9]|[2-9]\d|\d{3,})`
//	    Matches: 15, 16, 17, 18, 19, 20, 99, 100, 1000...
//	    Does not match: 0, 1, 10, 11, 12, 13, 14
//
//	NumGreaterOrEqual(100) → `(?:10[0-9]|1[1-9]\d|[2-9]\d{2}|\d{4,})`
//	    Matches: 100, 101, 150, 200, 999, 1000...
//	    Does not match: 0, 1, 50, 99
func NumGreaterOrEqual(n int) string {
	if n <= 0 {
		return VERSION_DIGITS
	}

	s := strconv.Itoa(n)
	numDigits := len(s)
	var patterns []string

	// Pattern for numbers with more digits than n (always greater)
	if numDigits < 10 {
		patterns = append(patterns, fmt.Sprintf(`\d{%d,}`, numDigits+1))
	}

	// Build patterns for numbers with the same number of digits as n
	for i := 0; i < numDigits; i++ {
		digit := int(s[i] - '0')
		prefix := s[:i]
		suffixLen := numDigits - i - 1

		if i == numDigits-1 {
			// Last digit: match this digit or higher
			patterns = append(patterns, prefix+digitRangeUp(digit))
		} else if digit < 9 {
			// Non-last digit: if higher, any following digits are fine
			patterns = append(patterns, prefix+digitRangeUp(digit+1)+fmt.Sprintf(`\d{%d}`, suffixLen))
		}
	}

	return joinPatterns(patterns)
}

// NumLessOrEqual generates a regex pattern that matches integers <= n.
//
// This function creates regex patterns for matching version number components
// that are less than or equal to a given value. It supports multi-digit
// numbers up to 10 digits (0 to 9999999999), making it suitable for semantic
// version comparisons where version parts can be arbitrarily large.
//
// Algorithm:
// The function builds multiple alternation patterns to cover all cases:
//  1. Numbers with fewer digits than n (always less, e.g., 0-9 when n=15)
//  2. Numbers with the same digit count but lower leading digits
//  3. Numbers matching the exact prefix with lower or equal subsequent digits
//
// Pattern Construction Strategy:
//   - For n<0: returns EMPTY_MATCH_PATTERN (matches nothing)
//   - For single digits (n=5): matches [0-5]
//   - For multi-digit (n=15): matches 0-9 (\d), 10-15 (1[0-5])
//   - For larger numbers: combines patterns for each digit position
//
// Parameters:
//   - n: The maximum value to match (inclusive). If n < 0, returns a pattern
//     that matches nothing.
//
// Returns:
//   - A regex pattern string that matches any integer <= n.
//     The pattern uses non-capturing groups (?:...) for alternation.
//
// Examples:
//
//	NumLessOrEqual(9)   → `[0-9]`
//	    Matches: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
//	    Does not match: 10, 11, 100...
//
//	NumLessOrEqual(15)  → `(?:\d|1[0-5])`
//	    Matches: 0, 1, 2, ..., 9, 10, 11, 12, 13, 14, 15
//	    Does not match: 16, 17, 20, 100...
//
//	NumLessOrEqual(123) → `(?:\d|\d{2}|1[0-1]\d|12[0-3])`
//	    Matches: 0, 1, 9, 10, 99, 100, 110, 119, 120, 121, 122, 123
//	    Does not match: 124, 125, 130, 200, 1000...
//
//	NumLessOrEqual(-1)  → `^(?!.*)$`
//	    Matches nothing (impossible constraint)
func NumLessOrEqual(n int) string {
	if n < 0 {
		return EMPTY_MATCH_PATTERN
	}

	s := strconv.Itoa(n)
	numDigits := len(s)
	var patterns []string

	// Pattern for numbers with fewer digits than n (always less)
	for d := 1; d < numDigits; d++ {
		if d == 1 {
			patterns = append(patterns, `\d`)
		} else {
			patterns = append(patterns, fmt.Sprintf(`\d{%d}`, d))
		}
	}

	// Build patterns for numbers with the same number of digits as n
	for i := 0; i < numDigits; i++ {
		digit := int(s[i] - '0')
		prefix := s[:i]
		suffixLen := numDigits - i - 1

		if i == numDigits-1 {
			// Last digit: match this digit or lower
			patterns = append(patterns, prefix+digitRangeDown(digit))
		} else if digit > 0 {
			// Non-last digit: if lower, any following digits are fine
			patterns = append(patterns, prefix+digitRangeDown(digit-1)+fmt.Sprintf(`\d{%d}`, suffixLen))
		}
	}

	return joinPatterns(patterns)
}

// digitRangeUp returns a pattern matching digits from d to 9.
func digitRangeUp(d int) string {
	if d == 9 {
		return "9"
	}
	return fmt.Sprintf("[%d-9]", d)
}

// digitRangeDown returns a pattern matching digits from 0 to d.
func digitRangeDown(d int) string {
	if d == 0 {
		return "0"
	}
	return fmt.Sprintf("[0-%d]", d)
}

// joinPatterns joins multiple patterns into a single regex with alternation.
func joinPatterns(patterns []string) string {
	if len(patterns) == 1 {
		return patterns[0]
	}
	return "(?:" + strings.Join(patterns, REGEX_OR) + ")"
}

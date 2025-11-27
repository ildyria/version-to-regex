package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// parseMavenRange parses Maven-style version ranges like [1.0,2.0), (1.0,2.0], etc.
func parseMavenRange(versionStr string) (*VersionConstraint, error) {
	// Maven ranges: [1.0,2.0), (1.0,2.0], [1.0,2.0], (1.0,2.0)
	// [ = inclusive lower bound, ( = exclusive lower bound
	// ] = inclusive upper bound, ) = exclusive upper bound

	if len(versionStr) < 3 {
		return nil, fmt.Errorf("invalid Maven range format: %s", versionStr)
	}

	openChar := versionStr[0]
	closeChar := versionStr[len(versionStr)-1]

	if (openChar != '[' && openChar != '(') || (closeChar != ']' && closeChar != ')') {
		return nil, fmt.Errorf("invalid Maven range brackets: %s", versionStr)
	}

	// Extract the version range content
	rangeContent := versionStr[1 : len(versionStr)-1]

	return &VersionConstraint{
		Operator: "maven-range",
		Version:  rangeContent,
	}, nil
}

// mavenRangeRegex creates a regex for Maven version ranges
func mavenRangeRegex(rangeStr string) (string, error) {
	lowerBound, upperBound, err := parseMavenRangeBounds(rangeStr)
	if err != nil {
		return "", err
	}

	switch {
	case lowerBound != "" && upperBound != "":
		return mavenBothBoundsPattern(lowerBound, upperBound), nil
	case lowerBound != "":
		return mavenLowerBoundPattern(lowerBound), nil
	case upperBound != "":
		return mavenUpperBoundPattern(upperBound), nil
	default:
		return EXACT_VERSION_TEMPLATE, nil
	}
}

// parseMavenRangeBounds extracts and validates the lower and upper bounds from a Maven range string
func parseMavenRangeBounds(rangeStr string) (lowerBound, upperBound string, err error) {
	parts := strings.Split(rangeStr, ",")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid Maven range format: %s", rangeStr)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

// extractMajorVersion extracts the major version number from a version string
func extractMajorVersion(version string) int {
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	return major
}

// mavenBothBoundsPattern creates a regex pattern when both bounds are specified
// Result (same major): ^1\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for [1.0,1.9])
// Result (diff major): ^[1-3]\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for [1.0,3.0])
func mavenBothBoundsPattern(lowerBound, upperBound string) string {
	lowerMajor := extractMajorVersion(lowerBound)
	upperMajor := extractMajorVersion(upperBound)

	if lowerMajor == upperMajor {
		return fmt.Sprintf(REGEX_START+"%d"+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS+VERSION_SUFFIX_PATTERN+REGEX_END, lowerMajor)
	}
	return fmt.Sprintf(REGEX_START+"[%d-%d]"+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS+VERSION_SUFFIX_PATTERN+REGEX_END, lowerMajor, upperMajor)
}

// mavenLowerBoundPattern creates a regex pattern when only lower bound is specified
// Result: ^(?:[2-9]|\d{2,})\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for [2.0,))
func mavenLowerBoundPattern(lowerBound string) string {
	lowerMajor := extractMajorVersion(lowerBound)
	return fmt.Sprintf(REGEX_START+`(?:[%d-9]|\d{2,})`+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS+VERSION_SUFFIX_PATTERN+REGEX_END, lowerMajor)
}

// mavenUpperBoundPattern creates a regex pattern when only upper bound is specified
// Result (major > 0): ^[0-3]\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for (,3.0])
// Result (major = 0): ^0\.\d+\.\d+(?:-[a-zA-Z0-9\-\.]+)?(?:\+[a-zA-Z0-9\-\.]+)?$ (for (,0.9])
func mavenUpperBoundPattern(upperBound string) string {
	upperMajor := extractMajorVersion(upperBound)
	if upperMajor > 0 {
		return fmt.Sprintf(REGEX_START+"[0-%d]"+VERSION_DOT+VERSION_DIGITS+VERSION_DOT+VERSION_DIGITS+VERSION_SUFFIX_PATTERN+REGEX_END, upperMajor)
	}
	return REGEX_START + "0" + VERSION_DOT + VERSION_DIGITS + VERSION_DOT + VERSION_DIGITS + VERSION_SUFFIX_PATTERN + REGEX_END
}

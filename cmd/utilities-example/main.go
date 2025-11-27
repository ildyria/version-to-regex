package main

import (
	"fmt"
	"log"

	"github.com/ildyria/version-to-regex/convert"
)

func main() {
	// Example 1: Basic usage with VersionMatches
	fmt.Println("=== VersionMatches Examples ===")

	versions := []string{"1.2.3", "1.2.5", "1.3.0", "2.0.0"}
	constraint := "^1.2.3"

	fmt.Printf("Testing versions against constraint '%s':\n", constraint)
	for _, version := range versions {
		matches, err := convert.VersionMatches(version, constraint)
		if err != nil {
			log.Printf("Error checking %s: %v", version, err)
			continue
		}
		status := "❌"
		if matches {
			status = "✅"
		}
		fmt.Printf("  %s %s\n", status, version)
	}

	// Example 2: Using MustVersionToRegex for compile-time constants
	fmt.Println("\n=== MustVersionToRegex Examples ===")

	// This would typically be done at package level as a constant
	var StableVersionPattern = convert.MustVersionToRegex("^1.0.0")

	testVersions := []string{"1.0.0", "1.2.3", "1.999.999", "2.0.0"}
	fmt.Println("Testing against stable version pattern (^1.0.0):")
	for _, version := range testVersions {
		matches := StableVersionPattern.MatchString(version)
		status := "❌"
		if matches {
			status = "✅"
		}
		fmt.Printf("  %s %s\n", status, version)
	}

	// Example 3: Different ecosystem constraints
	fmt.Println("\n=== Multi-ecosystem Examples ===")

	ecosystemExamples := []struct {
		ecosystem   string
		constraint  string
		version     string
		shouldMatch bool
	}{
		{"NPM", "^1.2.3", "1.5.0", true},
		{"NPM", "~1.2.3", "1.2.9", true},
		{"NPM", "~1.2.3", "1.3.0", false},
		{"Python", ">=1.2.3", "1.3.0", true},
		{"Python", "~=1.2.3", "1.2.9", true},
		{"Python", "~=1.2.3", "1.3.0", false},
		{"PHP/Composer", "^1.2.3", "1.9.0", true},
		{"PHP/Composer", "^1.2.3", "2.0.0", false},
	}

	for _, example := range ecosystemExamples {
		matches, err := convert.VersionMatches(example.version, example.constraint)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		status := "✅"
		if matches != example.shouldMatch {
			status = "❌"
		}

		fmt.Printf("%s: %s %s → %s (expected: %v, got: %v)\n",
			example.ecosystem,
			example.constraint,
			example.version,
			status,
			example.shouldMatch,
			matches)
	}
}

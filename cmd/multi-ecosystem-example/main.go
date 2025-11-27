package main

import (
	"fmt"
	"log"

	"github.com/ildyria/version-to-regex/convert"
)

func main() {
	fmt.Println("Multi-Ecosystem Version-to-Regex Examples")
	fmt.Println("=========================================")
	fmt.Println("")

	examples := []struct {
		ecosystem    string
		constraint   string
		testVersions []string
		description  string
	}{
		// NPM (Node.js)
		{
			ecosystem:    "NPM",
			constraint:   "^1.2.3",
			testVersions: []string{"1.2.3", "1.2.5", "1.3.0", "2.0.0"},
			description:  "NPM caret range (compatible within major version)",
		},
		{
			ecosystem:    "NPM",
			constraint:   "~1.2.3",
			testVersions: []string{"1.2.3", "1.2.5", "1.3.0"},
			description:  "NPM tilde range (compatible within minor version)",
		},

		// Python (pip)
		{
			ecosystem:    "Python",
			constraint:   ">=1.2.3",
			testVersions: []string{"1.2.2", "1.2.3", "1.3.0", "2.0.0"},
			description:  "Python greater than or equal",
		},
		{
			ecosystem:    "Python",
			constraint:   "~=1.2.3",
			testVersions: []string{"1.2.3", "1.2.5", "1.3.0"},
			description:  "Python compatible release",
		},

		// PHP (Composer)
		{
			ecosystem:    "PHP/Composer",
			constraint:   "^1.2.3",
			testVersions: []string{"1.2.3", "1.3.0", "1.999.999", "2.0.0"},
			description:  "Composer caret range",
		},

		// Maven (Java)
		{
			ecosystem:    "Maven",
			constraint:   "[1.0,2.0]",
			testVersions: []string{"0.9.0", "1.0.0", "1.5.0", "2.0.0", "3.0.0"},
			description:  "Maven version range (inclusive)",
		},
		{
			ecosystem:    "Maven",
			constraint:   "[1.0,]",
			testVersions: []string{"0.9.0", "1.0.0", "2.0.0", "3.0.0"},
			description:  "Maven lower bound only",
		},

		// Go Modules
		{
			ecosystem:    "Go",
			constraint:   "v1.2.3",
			testVersions: []string{"v1.2.3", "v1.2.3-beta.1", "v1.2.4", "1.2.3"},
			description:  "Go module semantic version",
		},
		{
			ecosystem:    "Go",
			constraint:   "v0.0.0-20210101000000-abcdef123456",
			testVersions: []string{"v0.0.0-20210101000000-abcdef123456", "v0.0.0-20210102000000-abcdef123456"},
			description:  "Go pseudo-version",
		},

		// C# NuGet
		{
			ecosystem:    "C#/NuGet",
			constraint:   "1.2.3.4567",
			testVersions: []string{"1.2.3.4567", "1.2.3.4568", "1.2.3"},
			description:  "C# 4-part version",
		},
		{
			ecosystem:    "C#/NuGet",
			constraint:   "1.0.0-alpha",
			testVersions: []string{"1.0.0-alpha", "1.0.0-beta", "1.0.0"},
			description:  "C# pre-release version",
		},

		// Wildcards (Universal)
		{
			ecosystem:    "Universal",
			constraint:   "1.2.*",
			testVersions: []string{"1.2.0", "1.2.999", "1.3.0"},
			description:  "Wildcard pattern",
		},
	}

	for _, example := range examples {
		fmt.Printf("### %s: %s\n", example.ecosystem, example.description)
		fmt.Printf("Constraint: %s\n", example.constraint)

		regex, err := convert.VersionToRegex(example.constraint)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Generated regex: %s\n", regex.String())
		fmt.Println("Test results:")

		for _, version := range example.testVersions {
			matches := regex.MatchString(version)
			status := "❌"
			if matches {
				status = "✅"
			}
			fmt.Printf("  %s %s\n", status, version)
		}

		fmt.Println("")
	}

	// Demonstrate the utility functions
	fmt.Println("=== Utility Functions Demo ===")

	// VersionMatches for quick checks
	matches, _ := convert.VersionMatches("v1.2.5", "v1.2.3")
	fmt.Printf("VersionMatches('v1.2.5', 'v1.2.3'): %v\n", matches)

	// MustVersionToRegex for compile-time constants
	var GoStablePattern = convert.MustVersionToRegex("v1.*")
	fmt.Printf("Go stable pattern matches v1.15.0: %v\n", GoStablePattern.MatchString("v1.15.0"))
	fmt.Printf("Go stable pattern matches v2.0.0: %v\n", GoStablePattern.MatchString("v2.0.0"))
}

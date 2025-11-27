package main

import (
	"fmt"
	"os"

	"github.com/ildyria/version-to-regex/convert"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: version-to-regex <version-constraint>")
		fmt.Println("Examples:")
		fmt.Println("  version-to-regex '>=1.2.3'")
		fmt.Println("  version-to-regex '^1.2.3'")
		fmt.Println("  version-to-regex '~1.2.3'")
		fmt.Println("  version-to-regex '1.*.0'")
		os.Exit(1)
	}

	constraint := os.Args[1]

	regex, err := convert.VersionToRegex(constraint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Version constraint: %s\n", constraint)
	fmt.Printf("Generated regex: %s\n", regex.String())

	// Test with example versions if provided
	if len(os.Args) > 2 {
		fmt.Println("\nTesting versions:")
		for _, version := range os.Args[2:] {
			matches := regex.MatchString(version)
			fmt.Printf("  %s: %v\n", version, matches)
		}
	}
}

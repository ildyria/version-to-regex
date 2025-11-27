# Version-to-regex

A comprehensive Go package that converts semantic version constraint strings to matching regular expressions. This package supports version constraint formats across multiple ecosystems including Python (pip), PHP (Composer), npm (Node.js), Maven (Java), Go modules, and C# NuGet.

## Features

- **Multiple constraint operators**: `=`, `==`, `>=`, `<=`, `>`, `<`, `!=`, `^`, `~`, `~>`, `~=`
- **Wildcard support**: `1.*`, `1.2.*`
- **NPM-style ranges**: Caret (`^`) and tilde (`~`) ranges
- **Python compatible release**: `~=` operator
- **Ruby pessimistic operator**: `~>` operator
- **Maven version ranges**: `[1.0,2.0]`, `(1.0,2.0)`, `[1.0,]`, `(,2.0]`
- **Go module versions**: `v1.2.3`, `v0.0.0-20210101000000-abcdef123456` (pseudo-versions)
- **C# NuGet versions**: `1.2.3.4567` (4-part), `1.0.0-alpha`, `1.0.0-preview`
- **Pre-release and build metadata support**: Handles `-alpha`, `+build` suffixes

## Installation

```bash
go get github.com/ildyria/version-to-regex
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "github.com/ildyria/version-to-regex/convert"
)

func main() {
    // Exact version match
    regex, err := convert.VersionToRegex("1.2.3")
    if err != nil {
        panic(err)
    }
    fmt.Println(regex.MatchString("1.2.3")) // true
    fmt.Println(regex.MatchString("1.2.4")) // false

    // NPM caret range (^1.2.3 allows 1.x.x but not 2.x.x)
    regex, err = convert.VersionToRegex("^1.2.3")
    if err != nil {
        panic(err)
    }
    fmt.Println(regex.MatchString("1.2.5")) // true
    fmt.Println(regex.MatchString("1.3.0")) // true
    fmt.Println(regex.MatchString("2.0.0")) // false

    // NPM tilde range (~1.2.3 allows 1.2.x but not 1.3.x)
    regex, err = convert.VersionToRegex("~1.2.3")
    if err != nil {
        panic(err)
    }
    fmt.Println(regex.MatchString("1.2.5")) // true
    fmt.Println(regex.MatchString("1.3.0")) // false
}
```

### As a CLI Tool

```bash
# Build the CLI tool
go build -o version-to-regex

# Test exact version
./version-to-regex "1.2.3" 1.2.3 1.2.4 1.3.0

# Test caret range
./version-to-regex "^1.2.3" 1.2.3 1.2.5 1.3.0 2.0.0

# Test tilde range
./version-to-regex "~1.2.3" 1.2.3 1.2.5 1.3.0

# Test wildcard
./version-to-regex "1.2.*" 1.2.0 1.2.999 1.3.0
```

## Supported Constraint Types

### Exact Match
- `1.2.3` or `==1.2.3` - Matches exactly version 1.2.3
- `=1.2.3` - Same as above

### Comparison Operators
- `>=1.2.3` - Greater than or equal to 1.2.3
- `<=1.2.3` - Less than or equal to 1.2.3
- `>1.2.3` - Greater than 1.2.3
- `<1.2.3` - Less than 1.2.3
- `!=1.2.3` - Not equal to 1.2.3

### Wildcard Patterns
- `1.*` - Any version with major version 1
- `1.2.*` - Any version with major.minor 1.2

### NPM-style Ranges
- `^1.2.3` - Compatible within the same major version (1.2.3 to < 2.0.0)
- `^0.2.3` - For 0.x versions, compatible within same minor (0.2.3 to < 0.3.0)
- `~1.2.3` - Compatible within the same minor version (1.2.3 to < 1.3.0)

### Python/Ruby Operators
- `~=1.2.3` - Python compatible release operator (same as tilde)
- `~>1.2.3` - Ruby pessimistic version operator (same as tilde)

## Examples by Ecosystem

### Python (pip)
```go
// Python-style constraints
convert.VersionToRegex(">=1.2.3")  // Greater than or equal
convert.VersionToRegex("~=1.2.3")  // Compatible release
convert.VersionToRegex("!=1.2.3")  // Not equal
```

### PHP (Composer)
```go
// Composer-style constraints
convert.VersionToRegex("^1.2.3")   // Caret operator
convert.VersionToRegex("~1.2.3")   // Tilde operator
convert.VersionToRegex(">=1.2.3")  // Greater than or equal
```

### Node.js (npm)
```go
// npm-style constraints
convert.VersionToRegex("^1.2.3")   // Caret range
convert.VersionToRegex("~1.2.3")   // Tilde range
convert.VersionToRegex("1.2.*")    // Wildcard
```

### Maven (Java)
```go
// Maven version ranges
convert.VersionToRegex("[1.0,2.0]")  // Inclusive range
convert.VersionToRegex("(1.0,2.0)")  // Exclusive range
convert.VersionToRegex("[1.0,]")     // Lower bound only
convert.VersionToRegex("(,2.0]")     // Upper bound only
```

### Go Modules
```go
// Go module versions
convert.VersionToRegex("v1.2.3")                              // Semantic version
convert.VersionToRegex("v0.0.0-20210101000000-abcdef123456")  // Pseudo-version
convert.VersionToRegex("v1.*")                                // Wildcard
```

### C# NuGet
```go
// C# NuGet versions
convert.VersionToRegex("1.2.3.4567")     // 4-part version
convert.VersionToRegex("1.0.0-alpha")    // Pre-release
convert.VersionToRegex("1.0.0-preview")  // Preview release
convert.VersionToRegex("1.0.0-rc.1")     // Release candidate
```

## Build and Test

```bash
# Run tests
go test ./...

# Build the CLI tool
go build -o version-to-regex

# Run with make (if Makefile is configured)
make build
make test
```

## API Reference

### `VersionToRegex(versionStr string) (*regexp.Regexp, error)`

Converts a semantic version constraint string to a compiled regular expression.

**Parameters:**
- `versionStr`: The version constraint string (e.g., "^1.2.3", ">=1.0.0")

**Returns:**
- `*regexp.Regexp`: Compiled regular expression that matches valid versions
- `error`: Error if the version string is invalid

### `VersionConstraint` struct

Represents a parsed version constraint with an operator and version.

```go
type VersionConstraint struct {
    Operator string  // The constraint operator (e.g., "^", ">=", "~")
    Version  string  // The version string (e.g., "1.2.3")
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the terms specified in your organization's license policy.
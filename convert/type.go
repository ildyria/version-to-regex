// Package convert provides version constraint types and utilities for converting
// semantic version strings to regular expression patterns.
//
// This file defines the core data structures used throughout the version
// conversion system to represent parsed version constraints from various
// package management ecosystems.
package convert

// Operator constants for version constraints
const (
	// OP_GREATER_EQUAL represents the >= operator (greater than or equal)
	OP_GREATER_EQUAL = ">="
	// OP_LESS_EQUAL represents the <= operator (less than or equal)
	OP_LESS_EQUAL = "<="
	// OP_NOT_EQUAL represents the != operator (not equal)
	OP_NOT_EQUAL = "!="
	// OP_EQUAL_EQUAL represents the == operator (exact match)
	OP_EQUAL_EQUAL = "=="
	// OP_PESSIMISTIC represents the ~> operator (Ruby pessimistic)
	OP_PESSIMISTIC = "~>"
	// OP_COMPATIBLE represents the ~= operator (Python compatible release)
	OP_COMPATIBLE = "~="
	// OP_GREATER represents the > operator (greater than)
	OP_GREATER = ">"
	// OP_LESS represents the < operator (less than)
	OP_LESS = "<"
	// OP_EQUAL represents the = operator (exact match)
	OP_EQUAL = "="
	// OP_CARET represents the ^ operator (NPM caret range)
	OP_CARET = "^"
	// OP_TILDE represents the ~ operator (NPM tilde range)
	OP_TILDE = "~"
	// OP_MAVEN_RANGE represents Maven version range with brackets
	OP_MAVEN_RANGE = "maven-range"
)

// VersionConstraint represents a semantic version constraint parsed from a version string.
//
// A version constraint consists of an operator that defines the relationship
// (equality, range, compatibility, etc.) and a version string that specifies
// the target version or version range.
//
// This structure is used internally to represent parsed version constraints
// before they are converted to regular expression patterns. It supports
// constraints from multiple package management ecosystems including:
//
// - NPM/Node.js (^, ~, exact matches)
// - Maven (version ranges with brackets)
// - Go modules (v-prefixed versions)
// - C# NuGet (4-part versions, pre-release patterns)
// - Python (compatible releases with ~=)
// - Ruby (pessimistic operator ~>)
//
// The Operator field contains the constraint type, while the Version field
// contains the version string or range specification.
//
// Examples of supported constraints:
//   - Exact match: Operator="==", Version="1.2.3"
//   - NPM caret: Operator="^", Version="1.2.3"
//   - NPM tilde: Operator="~", Version="1.2.3"
//   - Greater than: Operator=">", Version="1.2.3"
//   - Maven range: Operator="maven-range", Version="1.0,2.0"
//   - Python compatible: Operator="~=", Version="1.2.3"
//   - Ruby pessimistic: Operator="~>", Version="1.2.3"
type VersionConstraint struct {
	// Operator specifies the type of version constraint.
	//
	// Supported operators include:
	//   - "==" or "=": Exact version match
	//   - ">=": Greater than or equal to
	//   - "<=": Less than or equal to
	//   - ">": Greater than (exclusive)
	//   - "<": Less than (exclusive)
	//   - "!=": Not equal to
	//   - "^": NPM caret range (compatible within major version)
	//   - "~": NPM tilde range (compatible within minor version)
	//   - "~>": Ruby pessimistic operator (compatible release)
	//   - "~=": Python compatible release operator
	//   - "maven-range": Maven version range with brackets
	Operator string

	// Version contains the version string or range specification.
	//
	// The format and interpretation of this field depends on the Operator:
	//   - For exact matches: semantic version string (e.g., "1.2.3")
	//   - For comparison operators: semantic version string (e.g., "1.2.3")
	//   - For NPM ranges: semantic version string (e.g., "1.2.3")
	//   - For Maven ranges: comma-separated range (e.g., "1.0,2.0")
	//   - For Go modules: v-prefixed version (e.g., "v1.2.3")
	//   - For C# NuGet: may include 4-part versions (e.g., "1.2.3.4567")
	//
	// The version string may include pre-release identifiers and build metadata
	// following semantic versioning conventions, depending on the ecosystem.
	Version string
}

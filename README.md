# enum

Yet another Golang enum implementation. Follows a pattern I developed in the course of work.

A separate enum type allows you to add separate "tester" functions to the actual value itself.

```go
package foobar

import (
	"github.com/Riven-Spell/enum"
)

// Define our enum in a non-exported type
type eTestInt struct {
	// Value type, then enum struct type.
	// We do this as a embedded value to expose the functions.
	enum.EnumImpl[TestInt, eTestInt]
}

// Export our enum
var ETestInt eTestInt

// Define our value, exported
type TestInt int

// Export a stringifier
func (t TestInt) String() string {
	return ETestInt.String(t)
}

// Export a convenience parser
func (t *TestInt) Parse(s string) (err error) {
	*t, err = ETestInt.Parse(s)
	return
}

// Define our values
func (eTestInt) Foo() TestInt {
	return 1
}

func (eTestInt) Bar() TestInt {
	return 2
}

func (eTestInt) Baz() TestInt {
	return 3
}
```
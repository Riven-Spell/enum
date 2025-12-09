# enum

Yet another Golang enum implementation. Follows a pattern I developed in the course of work.

## Using a standard `EnumImpl`

The two type parameters are the result type of the enum, and the parent struct itself.

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

## Using a `BitflagEnumImpl`

Both `BitflagEnumImpl` and `BitflagImpl` take their backing `uint` type, the resulting flag type itself, and the parent enum struct.

Not returning bitflags by literal value adds a little complexity, but allows us to export convenience functions with it.

If you need the byte-for-byte value for whatever reason, `BitflagImpl` exports `Value() Raw`.

```go
package foobar

import (
	"github.com/Riven-Spell/enum"
)

type eTestBitflag struct {
	enum.BitflagEnumImpl[uint16, TestBitflag, eTestBitflag]
}

func (e eTestBitflag) GetDefaultBitflagStringOptions() enum.BitflagStringOptions {
	return enum.BitflagStringOptions{
		Separator: internal.Ptr("|"),
	}
}

var ETestBitFlag = eTestBitflag{}

func (e eTestBitflag) Foo() TestBitflag {
	return e.FromRawValue(1)
}

func (e eTestBitflag) Bar() TestBitflag {
	return e.FromRawValue(1 << 1)
}

func (e eTestBitflag) Baz() TestBitflag {
	return e.FromRawValue(1 << 2)
}

type TestBitflag struct {
	enum.BitflagImpl[uint16, TestBitflag, eTestBitflag]
}

func (t TestBitflag) String() string {
	return ETestBitFlag.String(t)
}
```
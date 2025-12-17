package enum

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type eTestBitflag struct {
	BitflagEnumImpl[uint16, TestBitflag, eTestBitflag]
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
	BitflagImpl[uint16, TestBitflag, eTestBitflag]
}

func (t TestBitflag) String() string {
	return ETestBitFlag.String(t)
}

func TestBitflagImpl(t *testing.T) {
	a := assert.New(t)

	foo := ETestBitFlag.Foo()
	bar := ETestBitFlag.Bar()
	baz := ETestBitFlag.Baz()
	fooBaz := foo.Add(baz)

	a.Equal(uint16(0b101), fooBaz.Value()) // validate basic operations
	a.True(fooBaz.Contains(foo))
	a.True(fooBaz.Contains(baz))
	a.False(fooBaz.Contains(bar))
	a.False(fooBaz.Add(bar).Remove(baz).Contains(baz))
	a.True(fooBaz.Add(bar).Remove(baz).Contains(bar))

	parsed, err := ETestBitFlag.Parse(ETestBitFlag.String(fooBaz), true)
	a.NoError(err)
	a.Equal(fooBaz.Value(), parsed.Value())
	a.True(strings.Contains(fooBaz.String(), ","))
	a.True(strings.Contains(ETestBitFlag.String(fooBaz), ","))

	a.Equal(foo, foo)
	a.NotEqual(foo, bar)
	a.True(foo == foo)
	a.False(foo == bar)

	switch foo.Add(baz) {
	case baz:
		a.Fail("triggered baz")
	case foo:
		a.Fail("triggered foo")
	case bar:
		a.Fail("triggered bar")
	case foo.Add(baz):
	default:
		a.Fail("defaulted")
	}
}

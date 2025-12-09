package enum_test

import (
	"strings"
	"testing"

	"github.com/Riven-Spell/enum"
	"github.com/Riven-Spell/enum/internal"
	"github.com/stretchr/testify/assert"
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

	parsed, err := ETestBitFlag.Parse(ETestBitFlag.String(fooBaz))
	a.NoError(err)
	a.Equal(fooBaz.Value(), parsed.Value())
	a.True(strings.Contains(fooBaz.String(), "|"))
}

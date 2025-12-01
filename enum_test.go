package enum_test

import (
	"testing"

	"github.com/Riven-Spell/enum"
	"github.com/stretchr/testify/assert"
)

type eTestInt struct {
	enum.EnumImpl[TestInt, eTestInt]
}

var ETestInt eTestInt

type TestInt int

func (t TestInt) String() string {
	return ETestInt.String(t)
}

func (t *TestInt) Parse(s string) (err error) {
	*t, err = ETestInt.Parse(s)
	return
}

func (eTestInt) Foo() TestInt {
	return 1
}

func (eTestInt) Bar() TestInt {
	return 2
}

func (eTestInt) Baz() TestInt {
	return 3
}

func TestEnumInt(t *testing.T) {
	a := assert.New(t)

	a.Equal(ETestInt.Foo().String(), "Foo")
	a.Equal(ETestInt.Bar().String(), "Bar")
	a.Equal(ETestInt.Baz().String(), "Baz")

	res, err := ETestInt.Parse("string")
	a.Error(err, "string should not parse")
	a.Equal(res, TestInt(0))

	res, err = ETestInt.Parse("foo")
	a.NoError(err, "foo should parse")
	a.Equal(res, TestInt(1))

	res, err = ETestInt.Parse("Bar")
	a.NoError(err, "Bar should parse")
	a.Equal(res, TestInt(2))

	res, err = ETestInt.Parse("baZ")
	a.NoError(err, "baZ should parse")
	a.Equal(res, TestInt(3))
}

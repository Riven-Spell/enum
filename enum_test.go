package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type eTestInt struct {
	EnumImpl[TestInt, eTestInt]
}

var ETestInt eTestInt

type TestInt int

func (t TestInt) String() string {
	return ETestInt.String(t)
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

	res, err := ETestInt.Parse("string", true)
	a.Error(err, "string should not parse")
	a.Equal(res, TestInt(0))

	res, err = ETestInt.Parse("foo", true)
	a.NoError(err, "foo should parse")
	a.Equal(res, TestInt(1))

	res, err = ETestInt.Parse("Bar", true)
	a.NoError(err, "Bar should parse")
	a.Equal(res, TestInt(2))

	res, err = ETestInt.Parse("baZ", true)
	a.NoError(err, "baZ should parse")
	a.Equal(res, TestInt(3))
}

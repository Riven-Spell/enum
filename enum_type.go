package enum

import (
	"fmt"
	"github.com/Riven-Spell/generic/enumerable"
	"strings"
)

// EnumVal adds a requirement that all values must be stringable.
// Implement it by calling back to the EnumImpl.String(Val) function.
type EnumVal interface {
	comparable
	fmt.Stringer
}

// EnumImpl should be embedded in a struct to treat it as an enum.
// Append functions to the pointer (or copied) value of the struct returning Vals
// in order to use.
type EnumImpl[Val EnumVal, Enum any] struct {
	valueNameCache map[Val]string
	nameValueCache map[string]Val
	typeName       string
}

func (e *EnumImpl[Val, Enum]) generateCaches() {
	globalRwLock.RLock() // grab a read lock, guarantee that the caches exist
	if e.nameValueCache != nil && e.valueNameCache != nil {
		globalRwLock.RUnlock()
		return
	}
	globalRwLock.RUnlock()

	// if they do not exist, grab the write lock and create them
	globalRwLock.Lock()
	defer globalRwLock.Unlock()

	e.nameValueCache, e.valueNameCache = generateCaches[Enum, Val, Val](noTransmute)
}

// String stringifies the input value.
func (e *EnumImpl[Val, Enum]) String(t Val) string {
	e.generateCaches()

	// Check the name value cache-- if it's a match, it will return the name of the function
	// if not, it'll return nothing.
	return e.valueNameCache[t]
}

// Parse parses a string to the resulting enum value.
func (e *EnumImpl[Val, Enum]) Parse(s string, strict bool) (v Val, err error) {
	e.generateCaches()

	// Lowercase our input first
	lowerInput := strings.ToLower(s)

	// Compare against the name value cache
	var ok bool
	v, ok = e.nameValueCache[lowerInput]
	if !ok && strict {
		err = fmt.Errorf("could not associate input `%s` with a value", s)
		return
	}

	return
}

func (e *EnumImpl[Val, Enum]) Values() []Val {
	e.generateCaches()

	return enumerable.Collect(enumerable.Map(enumerable.FromMap(e.valueNameCache), func(i enumerable.AB[Val, string]) Val {
		return i.A
	}))
}

package enum

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// EnumImpl should be embedded in a struct to treat it as an
type EnumImpl[Val comparable, Enum any] struct {
	cacheWriteOnce sync.Once

	valueNameCache map[Val]string
	nameValueCache map[string]Val
	typeName       string
}

func (e *EnumImpl[Val, Enum]) generateCaches() {
	e.cacheWriteOnce.Do(func() {
		// Make our maps first
		e.valueNameCache = make(map[Val]string)
		e.nameValueCache = make(map[string]Val)

		// Collect reflected types of our enum and value, enumRaw and vValue
		var enumRaw Enum
		var valueRaw Val
		eVal := reflect.ValueOf(enumRaw)
		eType := reflect.TypeOf(enumRaw)
		vType := reflect.TypeOf(valueRaw)

		// Step through the available methods on our enumeration type,
		// and find all that match our target function signature
		nMethods := eType.NumMethod()
		for i := 0; i < nMethods; i++ {
			// Get method[i] and it's type
			method := eType.Method(i)
			t := method.Type

			// Match our signature
			if !(t.NumIn() == 1 && t.NumOut() == 1 &&
				t.In(0).AssignableTo(eType) && t.Out(0).AssignableTo(vType)) {
				continue // One in (enum), one out (value
			}

			// put the name into our maps
			n := method.Name
			result := method.Func.Call([]reflect.Value{eVal})[0].Interface().(Val)
			e.valueNameCache[result] = n
			e.nameValueCache[strings.ToLower(n)] = result
		}
	})
}

func (e *EnumImpl[Val, Enum]) String(t Val) string {
	e.generateCaches()

	// Check the name value cache-- if it's a match, it will return the name of the function
	// if not, it'll return nothing.
	return e.valueNameCache[t]
}

func (e *EnumImpl[Val, Enum]) Parse(s string) (v Val, err error) {
	e.generateCaches()

	// Lowercase our input first
	lowerInput := strings.ToLower(s)

	// Compare against the name value cache
	var ok bool
	v, ok = e.nameValueCache[lowerInput]
	if !ok {
		err = fmt.Errorf("could not associate input `%s` with a value", s)
		return
	}

	return
}

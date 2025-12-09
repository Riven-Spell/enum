package enum

import (
	"reflect"
	"strings"
)

func noTransmute[T comparable](in T) T {
	return in
}

func generateCaches[Enum any, Interim any, Result comparable](transmuter func(Interim) Result) (nvCache map[string]Result, vnCache map[Result]string) {
	// Make our maps first
	vnCache = make(map[Result]string)
	nvCache = make(map[string]Result)

	// Collect reflected types of our enum and value, enumRaw and vValue
	var enum Enum
	var enumPtr = &enum
	var valueRaw Interim
	vType := reflect.TypeOf(valueRaw)

	type todo struct {
		Target reflect.Value
		Type   reflect.Type
	}

	todos := []todo{
		{reflect.ValueOf(enum), reflect.TypeOf(enum)},
		{reflect.ValueOf(enumPtr), reflect.TypeOf(enumPtr)},
	}

	for _, v := range todos {
		// Step through the available methods on our enumeration type,
		// and find all that match our target function signature
		nMethods := v.Type.NumMethod()
		for i := 0; i < nMethods; i++ {
			// Get method[i] and it's type
			method := v.Type.Method(i)
			t := method.Type

			// Match our signature
			if !(t.NumIn() == 1 && t.NumOut() == 1 &&
				t.In(0).AssignableTo(v.Type) && t.Out(0).AssignableTo(vType)) {
				continue // One in (enum), one out (value
			}

			// put the name into our maps
			n := method.Name
			interim := method.Func.Call([]reflect.Value{v.Target})[0].Interface().(Interim)
			result := transmuter(interim)
			vnCache[result] = n
			nvCache[strings.ToLower(n)] = result
		}
	}

	return
}

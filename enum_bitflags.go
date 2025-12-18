package enum

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Riven-Spell/enum/v2/internal"
)

type genericBfEnumImpl interface {
	bfEnumImpl()
}

// BitflagEnumImpl implements the bitflag enumerator. Paired with BitflagImpl, it allows for a similar structure to EnumImpl.
// To implement it, supply a uint "backing" type for the enum, a BitflagImpl "result" type, and the type encapsulating BitflagEnumImpl.
// BitflagImpl is required to add a fmt.Stringer to satisfy BitflagEnumImpl.
type BitflagEnumImpl[Raw internal.Flaggable, BfImpl genericBfImpl[Raw, Parent], Parent genericBfEnumImpl] struct {
	valueNameCache map[Raw]string
	nameValueCache map[string]Raw
	typeName       string
}

// FromRawValue returns the result type from a raw uint. This may not necessarily be a real value in the bitflags, so check your work!
func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) FromRawValue(val Raw) BfImpl {
	out := BitflagImpl[Raw, BfImpl, Parent]{}.getParentZeroInstance()
	ptr := getBitflagPtr(&out)

	ptr.value = val

	return out
}

func (BitflagEnumImpl[Raw, BfImpl, Parent]) bfEnumImpl() { panic("do not call, used for contracting") }

func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) generateCaches() {
	globalRwLock.RLock() // grab a read lock, guarantee that the caches exist
	if e.nameValueCache != nil && e.valueNameCache != nil {
		globalRwLock.RUnlock()
		return
	}
	globalRwLock.RUnlock()

	// if they do not exist, grab the write lock and create them
	globalRwLock.Lock()
	defer globalRwLock.Unlock()

	e.nameValueCache, e.valueNameCache = generateCaches[Parent, BfImpl, Raw](func(impl BfImpl) Raw {
		return impl.Value()
	})

	if noneName, ok := e.valueNameCache[0]; !ok {
		noneName = "None"
		e.nameValueCache[noneName] = 0
		e.valueNameCache[0] = noneName
	}
}

// BitflagSeparatorString defines the default separator string.
// Changing this will only lead to headaches, so it is a const.
const BitflagSeparatorString = ","

// String stringifies the target result type, using options set by ConfigureBitflagStringOptions.
func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) String(t BfImpl) string {
	e.generateCaches()

	results := make([]string, 0)
	for val, name := range e.valueNameCache {
		if val == 0 {
			continue // don't stringify zero until the end
		}

		if (t.Value() & val) == val {
			results = append(results, name)
		}
	}

	if len(results) == 0 { // If we have nothing, insert the zero value's name.
		results = append(results, e.valueNameCache[0])
	}

	return strings.Join(results, BitflagSeparatorString)
}

// Parse parses a string to the target result type. If options are not provided, falls back to two options:
// First, pulling them from DefaultBitflagStringOptionsGetter (if implemented on the parent struct)
// Second, GlobalDefaultBitflagStringOptions.
func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) Parse(s string, strict bool) (v BfImpl, err error) {
	e.generateCaches()

	v = BitflagImpl[Raw, BfImpl, Parent]{}.getParentZeroInstance()
	bfPtr := getBitflagPtr(&v)

	entriesRaw := strings.Split(s, BitflagSeparatorString)
	for _, raw := range entriesRaw {
		raw = strings.TrimSpace(raw)
		raw = strings.ToLower(raw)

		toAdd, ok := e.nameValueCache[raw]
		if !ok && strict {
			err = fmt.Errorf("could not associate input `%s` with a value", s)
			return
		}

		bfPtr.value |= toAdd
	}

	return
}

func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) Split(in BfImpl) []BfImpl {
	e.generateCaches()

	out := make([]BfImpl, 0)
	inVal := in.Value()

	for val, _ := range e.valueNameCache {
		if inVal&val == val {
			var toAdd BfImpl
			getBitflagPtr(&toAdd).value = val
			out = append(out, toAdd)
		}
	}

	return out
}

type genericBfImpl[F internal.Flaggable, E genericBfEnumImpl] interface {
	// contractual requirements, used for type inference
	bfImpl()
	enum() E

	Value() F

	fmt.Stringer
}

func getBitflagPtr[F internal.Flaggable, Enum genericBfEnumImpl, Parent genericBfImpl[F, Enum]](tgt *Parent) (out *BitflagImpl[F, Parent, Enum]) {
	derefType := reflect.TypeOf(out).Elem()

	parentPtrVal := reflect.ValueOf(tgt).Elem()
	nField := parentPtrVal.NumField()
	for i := 0; i < nField; i++ {
		f := parentPtrVal.Field(i)
		if f.Type().AssignableTo(derefType) {
			reflect.ValueOf(&out).Elem().Set(f.Addr())
			return out
		}
	}

	panic("could not find viable bitflag pointer (is BitflagImpl at the root of your struct?)")
}

// BitflagImpl is the companion type to BitflagEnumImpl. Both should be implemented together. Parameterization is the same.
// BitflagImpl is required to add a fmt.Stringer to satisfy BitflagEnumImpl.
type BitflagImpl[F internal.Flaggable, Parent genericBfImpl[F, Enum], Enum genericBfEnumImpl] struct {
	value F
}

func (b BitflagImpl[F, P, E]) bfImpl() { panic("do not call, used for contracting") }

func (b BitflagImpl[F, P, E]) enum() E { panic("do not call, used for type inference") }

func (b BitflagImpl[F, P, E]) getParentZeroInstance() (ret P) {
	return
}

// Value returns the raw uint backing value.
func (b BitflagImpl[F, P, E]) Value() F {
	return b.value
}

// Add "adds" two or more bitflags together in a binary OR operation.
func (b BitflagImpl[F, P, E]) Add(in ...P) P {
	out := b.getParentZeroInstance()
	bfPtr := getBitflagPtr(&out)
	bfPtr.value = b.value

	for _, v := range in {
		bfPtr.value |= v.Value()
	}

	return out
}

// Remove removes bitflags []in from the LHS bitflag and returns it, via bitwise AND + XOR of in.
func (b BitflagImpl[F, P, E]) Remove(in ...P) P {
	out := b.getParentZeroInstance()
	bfPtr := getBitflagPtr(&out)
	bfPtr.value = b.value

	for _, v := range in {
		bfPtr.value &= ^v.Value()
	}

	return out
}

// Contains returns whether or not a value exists in this bitflag.
func (b BitflagImpl[F, P, E]) Contains(in ...P) bool {
	for _, v := range in {
		if (b.value)&v.Value() != v.Value() {
			return false
		}
	}

	return true
}

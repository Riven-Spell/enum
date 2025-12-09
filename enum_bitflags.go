package enum

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Riven-Spell/enum/internal"
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
}

// String stringifies the target result type. If options are not provided, falls back to two options:
// First, pulling them from DefaultBitflagParseOptionsGetter (if implemented on the parent struct)
// Second, GlobalDefaultBitflagParseOptions.
func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) String(t BfImpl, opts ...BitflagStringOptions) string {
	e.generateCaches()

	var pType Parent

	opt := internal.FirstOrZero(opts)
	opt.SetDefaults(pType)

	results := make([]string, 0)
	for val, name := range e.valueNameCache {
		if (t.Value() & val) == val {
			results = append(results, name)
		}
	}

	return strings.Join(results, *opt.Separator)
}

// Parse parses a string to the target result type. If options are not provided, falls back to two options:
// First, pulling them from DefaultBitflagParseOptionsGetter (if implemented on the parent struct)
// Second, GlobalDefaultBitflagParseOptions.
func (e *BitflagEnumImpl[Raw, BfImpl, Parent]) Parse(s string, opts ...BitflagStringOptions) (v BfImpl, err error) {
	e.generateCaches()

	var pType Parent

	opt := internal.FirstOrZero(opts)
	opt.SetDefaults(pType)

	v = BitflagImpl[Raw, BfImpl, Parent]{}.getParentZeroInstance()
	bfPtr := getBitflagPtr(&v)

	entriesRaw := strings.Split(s, *opt.Separator)
	for _, raw := range entriesRaw {
		raw = strings.TrimSpace(raw)
		raw = strings.ToLower(raw)

		toAdd, ok := e.nameValueCache[raw]
		if !ok {
			err = fmt.Errorf("could not associate input `%s` with a value", s)
			return
		}

		bfPtr.value |= toAdd
	}

	return
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

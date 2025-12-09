package enum

import (
	"github.com/Riven-Spell/enum/internal"
)

// DefaultBitflagStringOptionsGetter can be implemented atop a struct including BitflagEnumImpl to provide default options for stringifying and parsing.
type DefaultBitflagStringOptionsGetter interface {
	GetDefaultBitflagStringOptions() BitflagStringOptions
}

// BitflagStringOptions specifies how bitflag strings should be formed.
type BitflagStringOptions struct {
	Separator *string
}

// GlobalDefaultBitflagStringOptions is the last fall-back option if no options are provided, and no default getter exists.
var GlobalDefaultBitflagStringOptions = BitflagStringOptions{
	Separator: internal.Ptr(","),
}

func (b *BitflagStringOptions) setDefaults(impl genericBfEnumImpl) {
	var defaults BitflagStringOptions
	if getter, ok := impl.(DefaultBitflagStringOptionsGetter); ok {
		defaults = getter.GetDefaultBitflagStringOptions()
	} else {
		defaults = GlobalDefaultBitflagStringOptions
	}

	if b.Separator == nil {
		b.Separator = defaults.Separator
	}
}

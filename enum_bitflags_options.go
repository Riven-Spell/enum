package enum

import (
	"github.com/Riven-Spell/enum/internal"
)

type DefaultBitflagStringOptionsGetter interface {
	GetDefaultBitflagStringOptions() BitflagStringOptions
}

type BitflagStringOptions struct {
	Separator *string
}

var defaultBitflagStringOptions = BitflagStringOptions{
	Separator: internal.Ptr(","),
}

func (b *BitflagStringOptions) SetDefaults(impl genericBfEnumImpl) {
	var defaults BitflagStringOptions
	if getter, ok := impl.(DefaultBitflagStringOptionsGetter); ok {
		defaults = getter.GetDefaultBitflagStringOptions()
	} else {
		defaults = defaultBitflagStringOptions
	}

	if b.Separator == nil {
		b.Separator = defaults.Separator
	}
}

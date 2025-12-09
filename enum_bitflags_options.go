package enum

import (
	"github.com/Riven-Spell/enum/internal"
)

type DefaultBitflagParseOptionsGetter interface {
	GetDefaultBitflagParseOptions() BitflagParseOptions
}

type BitflagParseOptions struct {
	Separator *string
}

var defaultBitflagParseOptions = BitflagParseOptions{
	Separator: internal.Ptr(","),
}

func (b *BitflagParseOptions) SetDefaults(impl genericBfEnumImpl) {
	var defaults BitflagParseOptions
	if getter, ok := impl.(DefaultBitflagParseOptionsGetter); ok {
		defaults = getter.GetDefaultBitflagParseOptions()
	} else {
		defaults = defaultBitflagParseOptions
	}

	if b.Separator == nil {
		b.Separator = defaults.Separator
	}
}

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

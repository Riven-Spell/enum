package enum

import (
	"errors"
	"strings"
	"sync"
)

/*
Why only allow a single set ever?

A) If caches are already generated, adjusting case sensitivity would require their regeneration.
B) It would be extremely annoying if a dependency were to change the semantics of enumeration.
C) It reduces any weird semantics of "how do we get options reliably?" to a contract set at library init.
*/

// BitflagStringOptions defines the set of options available for configuring bitflags
type BitflagStringOptions struct {
	// Separator determines the text that divides bitflag strings (i.e. Foo,Bar or Foo|Bar)
	Separator string
	// CaseInsensitive is relevant only in parsing; stringification will always result in the function name.
	CaseInsensitive bool
}

const (
	// DefaultBitflagSeparator defines the bitflag separator if ConfigureBitflagStringOptions is never called.
	DefaultBitflagSeparator = ","
	// DefaultBitflagCaseSens defines the case sensitivity (false = insensitive) if ConfigureBitflagStringOptions is never called.
	DefaultBitflagCaseSens = false
)

var bfStringOptionsConfigureOnce = &sync.Once{}
var bfStringOptions *BitflagStringOptions

var bfDefaultStringOptions = BitflagStringOptions{
	Separator:       DefaultBitflagSeparator,
	CaseInsensitive: DefaultBitflagCaseSens,
}

// ConfigureBitflagStringOptions permanently configures the bitflag string options.
// This should only ever be called
// a) intentionally, by the `main` package
// b) as a side effect of generateCaches.
// Middleware should never call this function.
func ConfigureBitflagStringOptions(cfg BitflagStringOptions) error {
	if strings.TrimSpace(cfg.Separator) == "" {
		return errors.New("separator cannot be empty")
	}

	bfStringOptionsConfigureOnce.Do(func() {
		bfStringOptions = &cfg
	})

	return nil
}

// GetBitflagStringOptions reads out the current bitflag string options.
func GetBitflagStringOptions() BitflagStringOptions {
	err := ConfigureBitflagStringOptions(bfDefaultStringOptions)
	if err != nil {
		panic("invalid default string options")
	}

	return *bfStringOptions
}

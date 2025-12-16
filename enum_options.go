package enum

import (
	"sync"
)

/*
Why only allow a single set ever?

A) If caches are already generated, adjusting case sensitivity would require their regeneration.
B) It would be extremely annoying if a dependency were to change the semantics of enumeration.
C) It reduces any weird semantics of "how do we get options reliably?" to a contract set at library init.
*/

// EnumStringOptions defines how enum field names get stringified and parsed.
type EnumStringOptions struct {
	// CaseInsensitive is relevant only in parsing; stringification will always result in the function name.
	CaseInsensitive bool
}

const (
	// DefaultEnumCaseSens case sensitive is off by default.
	DefaultEnumCaseSens bool = false
)

var enumStringOptions *EnumStringOptions
var enumStringOptionsSetOnce = &sync.Once{}
var enumDefaultStringOptions = EnumStringOptions{
	CaseInsensitive: DefaultEnumCaseSens,
}

// ConfigureEnumStringOptions permanently configures the bitflag string options.
// This should only ever be called
// a) intentionally, by the `main` package
// b) as a side effect of generateCaches.
// Middleware should never call this function.
func ConfigureEnumStringOptions(cfg EnumStringOptions) error {
	enumStringOptionsSetOnce.Do(func() {
		enumStringOptions = &cfg
	})

	// nothing is returned to avoid making this function a breaking change in the future.
	return nil
}

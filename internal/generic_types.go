package internal

// Flaggable encompasses unsigned integer types. Signed is intentionally not chosen.
type Flaggable interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

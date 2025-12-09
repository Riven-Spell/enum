package enum

import (
	"sync"
)

var globalRwLock = &sync.RWMutex{}

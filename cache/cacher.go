package cache

import "time"

// basic interface all our Caches have to implement
type Cacher interface{
	Set([]byte, []byte, time.Duration) error
	Get([]byte)([]byte, error)
	Has([]byte) bool
	Delete([]byte) error
}
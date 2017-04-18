// Package cow provides a reference implementation of copy-on-write maps
// for read-heavy, write-light data.
package cow

import (
	"sync"
	"sync/atomic"
)

type (
	Key   interface{}
	Value interface{}

	// Map is a synchronous copy on write map. Reads are cheap. Writes are expensive.
	Map struct {
		v atomic.Value
		m sync.Mutex // used only by writers
	}
)

// New initializes a new map based on an original map.
func New(m map[Key]Value) *Map {
	cowm := new(Map)
	cowm.v.Store(dup(m))
	return cowm
}

// NewEmpty initializes a new empty map.
func NewEmpty() *Map {
	return New(nil)
}

// Get retrieves the value associated with the key from the Map.
func (cowm *Map) Get(key Key) Value {
	m := cowm.v.Load().(map[Key]Value)
	return m[key]
}

// GetOK retrieves the value associated with the key from the Map.
func (cowm *Map) GetOK(key Key) (value Value, ok bool) {
	m := cowm.v.Load().(map[Key]Value)
	value, ok = m[key]
	return
}

// Set inserts a key-value pair.
func (cowm *Map) Set(key Key, value Value) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[Key]Value)
	dst := dup(src)
	dst[key] = value
	cowm.v.Store(dst)
}

// SetMap efficiently inserts all the values in m into the Map.
func (cowm *Map) SetMap(m map[Key]Value) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[Key]Value)
	dst := dup(src)
	copy(dst, m)
	cowm.v.Store(dst)
}

// Remove removes key from the Map.
func (cowm *Map) Remove(key Key) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[Key]Value)
	dst := dup(src)
	delete(dst, key)
	cowm.v.Store(dst)
}

// Clear removes all keys from the Map.
func (cowm *Map) Clear() {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	cowm.v.Store(make(map[Key]Value, 16))
}

// Reset initializes the Map to the values in m. Use of nil to empty the Map is okay.
func (cowm *Map) Reset(m map[Key]Value) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	cowm.v.Store(dup(m))
}

func copy(dst map[Key]Value, src map[Key]Value) {
	for k, v := range src {
		dst[k] = v
	}
}

func dup(src map[Key]Value) (dst map[Key]Value) {
	dst = make(map[Key]Value, len(src))
	copy(dst, src)
	return dst
}

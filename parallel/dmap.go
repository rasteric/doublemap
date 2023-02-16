// Package doublemap/parallel provides a generic parallel Map[K comparable, V comparable] with operations for getting and setting
// values by key, and the corresponding reverse map operation of getting and setting keys by values. The Map is
// thread-safe and uses an internal read/write mutex for synchronization. Otherwise the map works exactly the same as doublemap.
//
package parallel

import "sync"

type Map[K comparable, V comparable] struct {
	kv    map[K]V
	vk    map[V]K
	mutex sync.RWMutex
}

// New creates a new parallel double map.
func New[K, V comparable]() *Map[K, V] {
	return &Map[K, V]{kv: make(map[K]V), vk: make(map[V]K)}
}

// Get returns the value for the given key and true, the null value of the value type and false if no value
// was stored for this key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.kv[key]
	return value, ok
}

// Set sets a value for the given key.
func (m *Map[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.kv[key] = value
	m.vk[value] = key
}

// Remove removes the key and value mapping based on the given key. True is returned if the mapping was removed,
// false is returned when there was no mapping for the key in the first place.
func (m *Map[K, V]) Remove(key K) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, ok := m.Get(key)
	if ok {
		delete(m.kv, key)
		delete(m.vk, value)
		return true
	}
	return false
}

// ByValue returns the key for a given value and true, the key type's null value and false if no key was
// stored for this value.
func (m *Map[K, V]) ByValue(value V) (K, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	key, ok := m.vk[value]
	return key, ok
}

// RemoveByValue removes a given key-value mapping by the given value. True is returned if the mapping has been
// removed, false is returned if there was no such value in the double map in the first place.
func (m *Map[K, V]) RemoveByValue(value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	key, ok := m.ByValue(value)
	if ok {
		delete(m.kv, key)
		delete(m.vk, value)
		return true
	}
	return false
}

// Copy creates a copy of the key-value mapping. This operation is fairly slow but faster than using Get and Set
// manually. The copy is not deep, i.e., any key and values are just copied using ordinary assignment.
func (m *Map[K, V]) Copy() *Map[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	m2 := Map[K, V]{}
	for k, v := range m.kv {
		m2.kv[k] = v
		m2.vk[v] = k
	}
	return &m2
}

// Walk traverses key-value pairs in the map and provides them to the given function in unspecified order
// until the function returns false. The parallel map is read locked while walking it but not write locked.
func (m *Map[K, V]) Walk(fn func(key K, value V) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for k, v := range m.kv {
		if !fn(k, v) {
			break
		}
	}
}

// Clear clears the map, removing all key-valie pairs in it.
func (m *Map[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k := range m.kv { // better than one loop since this is optimized by compiler
		delete(m.kv, k)
	}
	for k := range m.vk {
		delete(m.vk, k)
	}
}

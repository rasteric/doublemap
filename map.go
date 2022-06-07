// Package doublemap provides a generic Map[K comparable, V comparable] with operations for getting and setting
// values by key, and the corresponding reverse map operation of getting and setting keys by values. The Map is
// not thread-safe.
//
// Example usage:
//     import "github.com/rasteric/doublemap"
//
//     func main() {
//       m := doublemap.Map[string]int
//       m.Set("hello", 1)
//       m.Set("world", 2)
//       m.SetByValue(3, "third")
//       m.Walk(func(key string, value int) bool {
//         fmt.Sprintf("%v, %v\n", key, value)
//         return true // will stop if return false
//       })
//     }
package doublemap

// A Map stores keys and values in a way that makes reverse mapping from values to keys efficient at the
// cost of additional memory and storage complexity. A Map does not have to be initialized, you can use it
// right out of the box. However, Map objects are not thread-safe.
type Map[K comparable,V comparable] struct {
	kv map[any]any
	vk map[any]any
}

func (m Map[K, V]) maybeInit() bool {
	if m.kv == nil {
		m.kv = make(map[any]any)
		m.vk = make(map[any]any)
		return true
	}
	return false
}

// Get returns the value for the given key and true, the null value of the value type and false if no value
// was stored for this key.
func (m Map[K, V]) Get(key K) (V, bool) {
  var result V
	if m.maybeInit() {
		return result, false
	}
	var ok bool
	var x any
	x, ok = m.kv[key]
	if !ok {
		return result, false
	}
	return x.(V), true
}

// Set sets a value for the given key. 
func (m Map[K, V]) Set(key K, value V) {
	m.maybeInit()
	m.kv[key] = value
	m.vk[value] = key
}

// Remove removes the key and value mapping based on the given key. True is returned if the mapping was removed,
// false is returned when there was no mapping for the key in the first place.
func (m Map[K, V]) Remove(key K) bool {
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
func (m Map[K, V]) ByValue(value V) (K, bool) {
	var result K
	if m.maybeInit() {
		return result, false
	}
	var ok bool
	var x any
	x, ok = m.vk[value]
	if !ok {
		return result, false
  }
	return x.(K), true
}

// RemoveByValue removes a given key-value mapping by the given value. True is returned if the mapping has been
// removed, false is returned if there was no such value in the double map in the first place.
func (m Map[K, V]) RemoveByValue(value V) bool {
	key, ok := m.ByValue(value)
	if ok {
		delete(m.kv, key)
		delete(m.vk, value)
		return true
	}
	return false
}

// SetByValue sets the key for the given value, i.e., it is like Set but in reverse direction.
func (m Map[K, V]) SetByValue(value V, key K) {
	oldkey, ok := m.ByValue(value)
	if ok {
		m.Remove(oldkey)
	}
	m.kv[key] = value
	m.vk[value] = key
}

// Copy creates a copy of the key-value mapping. This operation is fairly slow but faster than using Get and Set
// manually. The copy is not deep, i.e., any key and values are just copied using ordinary assignment.
func (m Map[K, V]) Copy() Map[K, V] {
	m2 := Map[K, V]{}
	m2.maybeInit()
	if m.maybeInit() {
    return m2
	}
	for k, v := range m.kv {
		m2.kv[k] = v
		m2.vk[v] = k 
	}
	return m2
}

// Walk traverses key-value pairs in the map and provides them to the given function in unspecified order
// until the function returns false. 
func (m Map[K, V]) Walk(fn func (key K, value V) bool) {
	for k, v := range m.kv {
		if !fn(k.(K), v.(V)) {
			break
		}
  }
}

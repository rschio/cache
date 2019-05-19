package cache

import "sync"

// Cache is a map[string]interface{}.
type Cache struct {
	sync.RWMutex
	m map[string]interface{}
	// max is the maximum number of elements
	// Cache can store.
	max int
}

// New creates a new Cache with maximum size max
// and at least size 1.
func New(max int) *Cache {
	if max <= 0 {
		max = 1
	}

	c := &Cache{}
	c.max = max
	c.m = make(map[string]interface{}, max)
	return c
}

// SetMax set a new max value.
func (c *Cache) SetMax(max int) {
	if max <= 0 {
		max = 1
	}
	c.Lock()
	c.max = max
	c.Unlock()

	c.resize()
}

func (c *Cache) resize() {
	c.Lock()
	defer c.Unlock()
	for k := range c.m {
		if len(c.m) <= c.max {
			break
		}
		delete(c.m, k)
	}
}

func (c *Cache) keepMaxSize() {
	c.Lock()
	defer c.Unlock()
	if len(c.m) < c.max {
		return
	}

	for k := range c.m {
		delete(c.m, k)
		break
	}
}

// Insert stores the element el with key k in the Cache. If the Cache
// is already full one element is deleted and then store el.
func (c *Cache) Insert(k string, el interface{}) {
	c.keepMaxSize()
	c.Lock()
	c.m[k] = el
	c.Unlock()
}

// Delete deletes the element of key k from Cache.
func (c *Cache) Delete(k string) {
	c.Lock()
	delete(c.m, k)
	c.Unlock()
}

// Get returns the element of key k and if it is in the Cache.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	el, ok := c.m[k]
	return el, ok
}

// Range range over the Cache and execute the function f
// for each element.
func (c *Cache) Range(f func(string, interface{})) {
	c.RLock()
	defer c.RUnlock()
	for key, val := range c.m {
		f(key, val)
	}
}

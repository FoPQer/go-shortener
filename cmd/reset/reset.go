package main

// Resetter is implemented by values that can restore their state to defaults.
type Resetter interface {
	Reset()
}

// Pool stores resettable items by key.
//
// T is the item type and must implement Resetter.
// K is the key type and must be comparable.
type Pool[T Resetter, K comparable] struct {
	items map[K]T
}

// NewPool creates an empty pool with initialized internal storage.
func NewPool[T Resetter, K comparable]() *Pool[T, K] {
	return &Pool[T, K]{
		items: make(map[K]T),
	}
}

// Get returns an item by key.
// If the key is missing, Get returns the zero value of T.
func (p *Pool[T, K]) Get(key K) T {
	if item, ok := p.items[key]; ok {
		return item
	}
	var zero T
	return zero
}

// Put stores or replaces an item by key.
func (p *Pool[T, K]) Put(key K, item T) {
	p.items[key] = item
}
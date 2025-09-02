package types

import (
	"slices"
)

type (
	OrderedMap[K comparable, V any] struct {
		insertionOrder []K
		values         map[K]V
		cmpFunc        func(a, b K) int
	}
)

// NewOrderedMap
//
// Initializes a new instance of the OrderedMap.
func NewOrderedMap[K comparable, V any](length int) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		insertionOrder: make([]K, 0, length),
		values:         make(map[K]V, length),
		cmpFunc:        nil,
	}
}

// NewOrderedMapWithCompareFunc
//
// Initializes a new instance of the OrderedMap with custom ordering
// dictated by the cmpFunc comparison function.
func NewOrderedMapWithCompareFunc[K comparable, V any](
	length int,
	cmpFunc func(a, b K) int,
) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		insertionOrder: make([]K, 0, length),
		values:         make(map[K]V, length),
		cmpFunc:        cmpFunc,
	}
}

// Len
//
// Returns the length of the internal map.
func (om *OrderedMap[K, V]) Len() int {
	return len(om.values)
}

// SetCompareFunc
//
// Sets the comparison function to determine ordering.
func (om *OrderedMap[K, V]) SetCompareFunc(cmpFunc func(a, b K) int) {
	om.cmpFunc = cmpFunc
}

// UnsetCompareFunc
//
// Unsets the comparison function.
func (om *OrderedMap[K, V]) UnsetCompareFunc() {
	om.cmpFunc = nil
}

// Get
//
// Gets the value associated with a key.
// If the key doesn't exist in the map the zero value
// for the type is returned and the boolean is set to false.
func (om OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := om.values[key]

	return v, ok
}

// Inserts a new value into the map while tracking the order.
func (om *OrderedMap[K, V]) Set(key K, value V) {
	_, exists := om.values[key]
	if !exists {
		om.insertionOrder = append(om.insertionOrder, key)
	}

	om.values[key] = value
}

// Unset
//
// Deletes the key from the map.
func (om *OrderedMap[K, V]) Unset(key K) {
	newInsertionOrder := make([]K, 0, len(om.insertionOrder)-1)
	for _, k := range om.insertionOrder {
		if k != key {
			newInsertionOrder = append(newInsertionOrder, k)
		}
	}

	om.insertionOrder = newInsertionOrder

	delete(om.values, key)
}

// Reset
//
// Resets all state including the comparison function.
func (om *OrderedMap[K, V]) Reset() {
	om.insertionOrder = make([]K, 0)
	om.values = make(map[K]V, 0)
	om.cmpFunc = nil
}

// Iter iterates through the map.
// If the CompareFunc is unset, it iterates in insertion order.
// Otherwise, it iterates in the order dictated by CompareFunc.
//
// Either use NewOrderedMapWithCompareFunc or
// call SetCompareFunc to set the comparison function.
func (om OrderedMap[K, V]) Iter(yield func(key K, value V) bool) {
	order := om.insertionOrder
	if om.cmpFunc != nil {
		order = make([]K, len(om.insertionOrder))
		copy(order, om.insertionOrder)

		slices.SortFunc(order, om.cmpFunc)
	}

	for _, key := range order {
		value := om.values[key]
		if !yield(key, value) {
			return
		}
	}
}

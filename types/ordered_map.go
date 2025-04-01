package types

import (
	"slices"
)

type (
	OrderedMap[T comparable, V any] struct {
		insertionOrder []T
		values         map[T]V
		cmpFunc        func(a, b T) int
	}
)

// NewOrderedMap
//
// Initializes a new instance of the OrderedMap
func NewOrderedMap[T comparable, V any](length int) OrderedMap[T, V] {
	return OrderedMap[T, V]{
		insertionOrder: make([]T, 0, length),
		values:         make(map[T]V, length),
		cmpFunc:        nil,
	}
}

func NewOrderedMapWithCompareFunc[T comparable, V any](
	length int,
	cmpFunc func(a, b T) int,
) OrderedMap[T, V] {
	return OrderedMap[T, V]{
		insertionOrder: make([]T, 0, length),
		values:         make(map[T]V, length),
		cmpFunc:        cmpFunc,
	}
}

// SetCompareFunc
//
// Sets the comparison function to determine ordering.
func (om *OrderedMap[T, V]) SetCompareFunc(cmpFunc func(a, b T) int) {
	om.cmpFunc = cmpFunc
}

// UnsetCompareFunc
//
// Unsets the comparison function
func (om *OrderedMap[T, V]) UnsetCompareFunc() {
	om.cmpFunc = nil
}

func (om OrderedMap[T, V]) Get(key T) (V, bool) {
	v, ok := om.values[key]

	return v, ok
}

func (om *OrderedMap[T, V]) Set(key T, value V) {
	_, exists := om.values[key]
	if !exists {
		om.insertionOrder = append(om.insertionOrder, key)
	}

	om.values[key] = value
}

func (om *OrderedMap[T, V]) Unset(key T) {
	newInsertionOrder := make([]T, 0, len(om.insertionOrder)-1)
	for _, k := range om.insertionOrder {
		if k != key {
			newInsertionOrder = append(newInsertionOrder, k)
		}
	}

	om.insertionOrder = newInsertionOrder

	delete(om.values, key)
}

func (om *OrderedMap[T, V]) Reset() {
	om.insertionOrder = make([]T, 0)
	om.values = make(map[T]V, 0)
	om.cmpFunc = nil
}

// Iter iterates through the map.
// If the CompareFunc is unset, it iterates in insertion order.
// Otherewise, it iterates in the order dictated by CompareFunc.
//
// Call SetCompareFunc to set the comparison function.
func (om OrderedMap[T, V]) Iter(yield func(key T, value V) bool) {
	order := om.insertionOrder
	if om.cmpFunc != nil {
		order = make([]T, len(om.insertionOrder))
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

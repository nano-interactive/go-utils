package types

type (
	OrderedMap[T comparable, V any] struct {
		insertionOrder []T
		values         map[T]V
	}
)

func NewOrderedMap[T comparable, V any](length int) OrderedMap[T, V] {
	return OrderedMap[T, V]{
		insertionOrder: make([]T, 0, length),
		values:         make(map[T]V, length),
	}
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
}

func (om OrderedMap[T, V]) Iter(yield func(key T, Value V) bool) {
	for _, key := range om.insertionOrder {
		value := om.values[key]
		if !yield(key, value) {
			return
		}
	}
}

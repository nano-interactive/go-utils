package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"slices"
)

var (
	ErrInvalidJSON    = errors.New("expected a JSON object")
	ErrInvalidJSONKey = errors.New("expected a string key in JSON")
)

// UnmarshalJSON
//
// Unserializes JSON objects into an ordered map
// preserving the order in which keys appear in
// the input datastream.
//
// Doesn't support unmarshaling JSON arrays as those
// are already ordered structures for which you
// shouldn't use this ordered map.
func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	delimiter, ok := token.(json.Delim)
	if !ok || delimiter != '{' {
		return ErrInvalidJSON
	}

	if om.insertionOrder == nil {
		om.insertionOrder = make([]K, 0)
	}

	if om.values == nil {
		om.values = make(map[K]V, 0)
	}

	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			return err
		}

		key, ok := token.(string)
		if !ok {
			return ErrInvalidJSONKey
		}

		kJson, err := json.Marshal(key)
		if err != nil {
			return err
		}

		var k K
		if err := json.Unmarshal(kJson, &k); err != nil {
			return err
		}

		var v V
		if err := decoder.Decode(&v); err != nil {
			return err
		}

		om.Set(k, v)
	}

	_, err = decoder.Token()
	return err
}

// MarshalJSON
//
// Serializes the ordered map into a JSON stream.
//
// If the CompareFunc is unset, it serializes in insertion order.
// Otherwise, it serializes in the order dictated by CompareFunc.
//
// Either use NewOrderedMapWithCompareFunc or
// call SetCompareFunc to set the comparison function.
func (om OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	if len(om.values) == 0 {
		return []byte("{}"), nil
	}

	b := new(bytes.Buffer)
	order := om.insertionOrder
	if om.cmpFunc != nil {
		order = make([]K, len(om.insertionOrder))
		copy(order, om.insertionOrder)

		slices.SortFunc(order, om.cmpFunc)
	}

	b.WriteRune('{')
	for i, k := range order {
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(om.values[k])
		if err != nil {
			return nil, err
		}

		b.Write(key)
		b.WriteRune(':')
		b.Write(value)
		if i < len(om.values)-1 {
			b.WriteRune(',')
		}
	}
	b.WriteRune('}')

	return b.Bytes(), nil
}

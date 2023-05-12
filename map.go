package sorted

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
)

// NewMap creates a new map.
func NewMap[TKey comparable, TValue any]() *Map[TKey, TValue] {
	return &Map[TKey, TValue]{
		mapKeysToIndex: make(map[TKey]*mapValue[TKey, TValue]),
	}
}

// Map is a map which retains the order of the keys.
type Map[TKey comparable, TValue any] struct {
	mapKeysToIndex map[TKey]*mapValue[TKey, TValue]
	index          int64
}

// Add an item to the map.
func (m *Map[TKey, TValue]) Add(k TKey, v TValue) {
	if kv, ok := m.mapKeysToIndex[k]; ok {
		kv.value = v
		return
	}
	m.mapKeysToIndex[k] = &mapValue[TKey, TValue]{
		index: atomic.AddInt64(&m.index, 1),
		key:   k,
		value: v,
	}
}

// Get an item from the map.
func (m *Map[TKey, TValue]) Get(k TKey) (v TValue, ok bool) {
	kv, ok := m.mapKeysToIndex[k]
	if !ok {
		return v, false
	}
	return kv.value, true
}

// Del deletes an item from the map.
func (m *Map[TKey, TValue]) Del(k TKey) {
	delete(m.mapKeysToIndex, k)
}

// Keys returns all of the keys within the map.
func (m *Map[TKey, TValue]) Keys() (keys []TKey) {
	kvs := make([]*mapValue[TKey, TValue], len(m.mapKeysToIndex))
	var index int
	for _, kv := range m.mapKeysToIndex {
		kvs[index] = kv
		index++
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].index < kvs[j].index
	})
	keys = make([]TKey, len(m.mapKeysToIndex))
	for i, v := range kvs {
		keys[i] = v.key
	}
	return keys
}

type mapValue[TKey comparable, TValue any] struct {
	index int64
	key   TKey
	value TValue
}

// MarshalJSON marshals to JSON.
func (m *Map[TKey, TValue]) MarshalJSON() (result []byte, err error) {
	var k TKey
	if reflect.TypeOf(k) != reflect.TypeOf("") {
		return nil, &json.UnsupportedTypeError{
			Type: reflect.TypeOf(k),
		}
	}
	var buf bytes.Buffer
	if _, err = buf.WriteString("{"); err != nil {
		return
	}
	keys := m.Keys()
	var kb []byte
	for i, k := range keys {
		// "key"
		if kb, err = json.Marshal(k); err != nil {
			return
		}
		if _, err = buf.Write(kb); err != nil {
			return
		}
		// :
		if _, err = buf.WriteString(":"); err != nil {
			return
		}
		// "value"
		v, _ := m.Get(k)
		if kb, err = json.Marshal(v); err != nil {
			return
		}
		if _, err = buf.Write(kb); err != nil {
			return
		}
		// Skip trailing comma.
		if i == len(keys)-1 {
			break
		}
		// ","
		if _, err = buf.WriteString(","); err != nil {
			return
		}
	}
	if _, err = buf.WriteString("}"); err != nil {
		return
	}
	return json.RawMessage(buf.Bytes()), nil
}

// UnmarshalJSON marshals from JSON.
func (m *Map[TKey, TValue]) UnmarshalJSON(data []byte) (err error) {
	d := json.NewDecoder(bytes.NewReader(data))
	// We expect...
	// {
	t, err := d.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return fmt.Errorf("expected map start delimiter ('{') not found, got %T", t)
	}
	if delim.String() != "{" {
		return fmt.Errorf("expected map start delimiter ('{') not found, got %q", delim.String())
	}

	nm := NewMap[TKey, TValue]()

	// Next, get keys and values, until we get a close brace delimiter.
	for {
		t, err = d.Token()
		if err != nil {
			return err
		}
		delim, ok = t.(json.Delim)
		if ok {
			if delim.String() != "}" {
				return fmt.Errorf("expected map close delimiter ('}') not found, got %q", delim.String())
			}
			break
		}
		// Read key and value.
		var k TKey = t.(TKey)
		var v TValue
		if err = d.Decode(&v); err != nil {
			return
		}
		nm.Add(k, v)
	}

	(*m) = (*nm)
	return
}

package sorted

import (
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

package sorted

import (
	"encoding/json"
	"sort"
	"sync/atomic"
)

// NewSet creates a new set.
func NewSet[T comparable](values ...T) *Set[T] {
	ss := &Set[T]{
		mapKeysToIndex: make(map[T]int64),
	}
	for _, v := range values {
		ss.Add(v)
	}
	return ss
}

// Set which retains the order that the keys were added.
type Set[T comparable] struct {
	mapKeysToIndex map[T]int64
	index          int64
}

// Add an item to the set.
func (s *Set[T]) Add(v T) {
	if _, ok := s.mapKeysToIndex[v]; ok {
		return
	}
	s.mapKeysToIndex[v] = atomic.AddInt64(&s.index, 1)
}

// Contains determines whether an item is in the set.
func (s *Set[T]) Contains(v T) (ok bool) {
	_, ok = s.mapKeysToIndex[v]
	return
}

// Del deletes an item from the set.
func (s *Set[T]) Del(v T) {
	delete(s.mapKeysToIndex, v)
}

// Values returns all of the values within the set.
func (s *Set[T]) Values() (v []T) {
	values := make([]setValue[T], len(s.mapKeysToIndex))
	var index int
	for k, v := range s.mapKeysToIndex {
		values[index] = setValue[T]{
			index: v,
			value: k,
		}
		index++
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].index < values[j].index
	})
	v = make([]T, len(s.mapKeysToIndex))
	for i, vv := range values {
		v[i] = vv.value
	}
	return
}

type setValue[T comparable] struct {
	index int64
	value T
}

// MarshalJSON marshals to JSON.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON marshals from JSON.
func (s *Set[T]) UnmarshalJSON(data []byte) (err error) {
	var a []T
	if err = json.Unmarshal(data, &a); err != nil {
		return err
	}
	(*s) = (*NewSet(a...))
	return
}

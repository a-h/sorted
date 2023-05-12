package sorted

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleMap() {
	m := NewMap[string, string]()
	m.Add("a", "aa")
	m.Add("b", "bb")
	m.Add("b", "bbb")
	m.Add("c", "c")
	m.Del("c")
	value, ok := m.Get("b")
	fmt.Println(value, ok, m.Keys())
	// Output: bbb true [a b]
}

func TestStringToStringMap(t *testing.T) {
	m := NewMap[string, string]()
	m.Add("a", "aa")
	m.Add("b", "bb")
	m.Add("c", "cc")
	if !reflect.DeepEqual(m.Keys(), []string{"a", "b", "c"}) {
		t.Errorf("Expected map to only contain a, b, c, got %v", m.Keys())
	}
	m.Add("a", "aaa")
	aaa, ok := m.Get("a")
	if !ok {
		t.Errorf("Expected to be able to get 'a', but couldn't")
	}
	if aaa != "aaa" {
		t.Errorf("Expected to get value 'aaa' but got '%v'", aaa)
	}
	m.Del("c")
	if !reflect.DeepEqual(m.Keys(), []string{"a", "b"}) {
		t.Errorf("Expected m to only contain a, b, got %v", m.Keys())
	}
	m.Add("c", "123")
	if !reflect.DeepEqual(m.Keys(), []string{"a", "b", "c"}) {
		t.Errorf("Expected m to contain a, b, c after restore, got %v", m.Keys())
	}
	c, ok := m.Get("c")
	if !ok {
		t.Errorf("Expected to be able to get the value, but couldn't")
	}
	if c != "123" {
		t.Errorf("Expected to be able to get the value, got %v", c)
	}
	m.Del("b")
	if !reflect.DeepEqual(m.Keys(), []string{"a", "c"}) {
		t.Errorf("Expected map to only contain a, c, got %v", m.Keys())
	}
	if _, ok := m.Get("d"); ok {
		t.Error("Expected not to be able to get value 'd' that doesn't exist, but returned true")
	}
}

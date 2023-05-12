package sorted

import (
	"encoding/json"
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

func TestMap(t *testing.T) {
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

func TestMapJSON(t *testing.T) {
	t.Run("marshalling maps is supported", func(t *testing.T) {
		m := NewMap[string, string]()
		m.Add("a", "aa")
		m.Add("b", "bb")
		m.Add("c", "cc")
		actual, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		if string(actual) != `{"a":"aa","b":"bb","c":"cc"}` {
			t.Errorf("unexpected JSON: %q", string(actual))
		}
	})
	t.Run("unmarshalling maps is supported", func(t *testing.T) {
		var m *Map[string, string]
		err := json.Unmarshal([]byte(`{"c":"cc","b":"bb","a":"aa"}`), &m)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}
		if !reflect.DeepEqual(m.Keys(), []string{"c", "b", "a"}) {
			t.Errorf("expected m to contain c, b, a, got %v", m.Keys())
		}
	})
	t.Run("unmarshalling invalid JSON results in an error", func(t *testing.T) {
		var m *Map[string, string]
		err := json.Unmarshal([]byte(`--...{`), &m)
		if err == nil {
			t.Errorf("expected JSON unmarshal error not returned")
		}
	})
	t.Run("marshalling null maps is supported", func(t *testing.T) {
		var m *Map[string, string]
		marshalled, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		if string(marshalled) != "null" {
			t.Fatalf("expected 'null', got %v", string(marshalled))
		}
	})
	t.Run("unmarshalling null maps is supported", func(t *testing.T) {
		var m *Map[string, string]
		err := json.Unmarshal([]byte("null"), &m)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if m != nil {
			t.Fatalf("expected nil, got %v", m)
		}
	})
	t.Run("marshalling maps with string keys is supported", func(t *testing.T) {
		m := NewMap[string, string]()
		m.Add("key1", "value1")
		m.Add("key2", "value2")
		marshalled, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		expected := `{"key1":"value1","key2":"value2"}`
		if string(marshalled) != expected {
			t.Fatalf("expected %v, got %v", expected, string(marshalled))
		}
	})
}

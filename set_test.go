package sorted

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func ExampleSet() {
	set := NewSet("a", "b")
	set.Add("d")
	set.Add("c")
	set.Del("d")
	fmt.Println(set.Values())
	// Output: [a b c]
}

func TestSet(t *testing.T) {
	set := NewSet("a", "b", "c")
	if !reflect.DeepEqual(set.Values(), []string{"a", "b", "c"}) {
		t.Errorf("Expected set to only contain a, b, c, got %v", set.Values())
	}
	set.Add("a")
	ok := set.Contains("a")
	if !ok {
		t.Errorf("Expected to be able to get the value 'a', but couldn't")
	}
	set.Del("c")
	if !reflect.DeepEqual(set.Values(), []string{"a", "b"}) {
		t.Errorf("Expected set to only contain a, b, got %v", set.Values())
	}
	set.Add("c")
	if !reflect.DeepEqual(set.Values(), []string{"a", "b", "c"}) {
		t.Errorf("Expected set to contain a, b, c after restore, got %v", set.Values())
	}
	set.Del("b")
	if !reflect.DeepEqual(set.Values(), []string{"a", "c"}) {
		t.Errorf("Expected set to only contain a, c, got %v", set.Values())
	}
}

func TestSetJSON(t *testing.T) {
	t.Run("marshalling sets is supported", func(t *testing.T) {
		set := NewSet("a", "b", "c")
		actual, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		if string(actual) != `["a","b","c"]` {
			t.Errorf("unexpected JSON: %q", string(actual))
		}
	})
	t.Run("unmarshalling sets is supported", func(t *testing.T) {
		var set *Set[string]
		err := json.Unmarshal([]byte(`["d", "c", "a"]`), &set)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}
		if !reflect.DeepEqual(set.Values(), []string{"d", "c", "a"}) {
			t.Errorf("unexpected JSON unmarshal result: %#v", set.Values())
		}
	})
	t.Run("unmarshalling invalid JSON results in an error", func(t *testing.T) {
		var set *Set[string]
		err := json.Unmarshal([]byte(`--...{`), &set)
		if err == nil {
			t.Errorf("expected JSON unmarshal error not returned")
		}
	})
	t.Run("marshalling null sets is supported", func(t *testing.T) {
		var set *Set[string]
		marshalled, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		if string(marshalled) != "null" {
			t.Fatalf("expected 'null', got %v", string(marshalled))
		}
	})
	t.Run("unmarshalling null sets is supported", func(t *testing.T) {
		var set *Set[string]
		err := json.Unmarshal([]byte("null"), &set)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if set != nil {
			t.Fatalf("expected nil, got %v", set)
		}
	})
}

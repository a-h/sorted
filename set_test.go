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
	set := NewSet("a", "b", "c")
	actual, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	if string(actual) != `["a","b","c"]` {
		t.Errorf("Unexpected JSON: %q", string(actual))
	}
	err = json.Unmarshal([]byte(`["d", "c", "a"]`), &set)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if !reflect.DeepEqual(set.Values(), []string{"d", "c", "a"}) {
		t.Errorf("Unexpected JSON unmarshal result: %#v", set.Values())
	}
	err = json.Unmarshal([]byte(`--...{`), &set)
	if err == nil {
		t.Errorf("Expected JSON unmarshal error not returned")
	}
}

# sorted

Set and map which retains the order of the inserted keys.

## Map

```go
m := NewMap[string, string]()
m.Add("a", "aa")
m.Add("b", "bb")
m.Add("b", "bbb")
m.Add("c", "c")
m.Del("c")
value, ok := m.Get("b")
fmt.Println(value, ok, m.Keys())
// Output: bbb true [a b]
```

## Set

```go
set := NewSet("a", "b")
set.Add("d")
set.Add("c")
set.Del("d")
fmt.Println(set.Values())
// Output: [a b c]
```


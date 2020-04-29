package pj

import (
	"testing"
)

func TestParseSimpleJSON(t *testing.T) {
	testcases := []struct {
		input             string
		expectKey         string
		expectVal         interface{}
		expectChildrenLen int
	}{
		{`{ "item1": 42 }`, "item1", 42, 1},
		{`{ "item2": "test2" }`, "item2", "test2", 1},
		{`{ "item3": true }`, "item3", true, 1},
		{`{ "item4": false }`, "item4", false, 1},
	}

	for _, tt := range testcases {
		l := NewLexer(tt.input)
		p := NewParser(l)

		root, err := p.Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v\n", err)
		}

		checkParseSimpleJSON(t, root, tt.expectKey, tt.expectVal, tt.expectChildrenLen)
	}
}

func checkParseSimpleJSON(t *testing.T, r *Root, k string, v interface{}, l int) {
	rootValue := *r.val
	switch rootValue.(type) {
	case Object:
		obj, _ := rootValue.(Object)
		if l != len(obj.children) {
			t.Fatalf("Object Children Length: expect is %d but got %d\n", l, len(obj.children))
		}

		for _, c := range obj.children {
			if c.key != k {
				t.Fatalf("Object Children Key: expect is %s but got %s\n", k, c.key)
			}
			if c.val != v {
				t.Fatalf("Object Children Value: expect is %s but got %s\n", v, c.val)
			}
		}
	}
}

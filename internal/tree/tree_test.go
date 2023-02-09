package tree

import (
	"reflect"
	"testing"
)

func TestTree(t *testing.T) {
	separators := []byte{'/'}
	t1 := New("", separators)
	t1.AddKey("foo/bar", 1)
	t1.AddKey("foo/bar/baz", 2)
	t1.AddKey("foo/qux", 3)
	if size := t1.GetSize("foo/"); size != 6 {
		t.Errorf("Expected size of 6, got %d", size)
	}
	if size := t1.GetSize("foo/bar/"); size != 2 {
		t.Errorf("Expected size of 2, got %d", size)
	}
	expanded := t1.Expand("foo/")
	expected := map[string]int64{
		"bar":  1,
		"bar/": 2,
		"qux":  3,
	}
	if !reflect.DeepEqual(expanded, expected) {
		t.Errorf("Expected expanded to be %v, got %v", expected, expanded)
	}
}

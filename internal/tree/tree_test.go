package tree

import (
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
}

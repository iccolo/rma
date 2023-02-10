package tree

import (
	"fmt"
	"sort"
)

func New(name string, separators []byte) *Tree {
	t := &Tree{
		root: &Node{
			Segment: name,
			Child:   make(map[string]*Node),
		},
		separators: make(map[byte]bool),
	}
	for _, separator := range separators {
		t.separators[separator] = true
	}
	return t
}

type Tree struct {
	root       *Node
	nodeNum    int64
	separators map[byte]bool
}

type Node struct {
	Segment string           // key segment
	KeyNum  int64            // child key num
	Size    int64            // total size of keys in current and child node
	Child   map[string]*Node // child node
}

func (t *Tree) AddKey(key string, size int64) {
	tmpRoot := t.root
	left := 0
	right := left
	for right < len(key) {
		// key end
		if right == len(key)-1 {
			segment := key[left : right+1]
			if child, ok := tmpRoot.Child[segment]; ok { // exists duplicate key, cover duplicate key size
				tmpRoot.Size += size - child.Size
				child.Size = size
			} else {
				if tmpRoot.Child == nil {
					tmpRoot.Child = make(map[string]*Node)
				}
				newNode := &Node{
					Segment: segment,
					Size:    size,
					KeyNum:  1,
				}
				tmpRoot.Child[segment] = newNode
				tmpRoot.Size += size
				tmpRoot.KeyNum++
				tmpRoot = newNode
				t.nodeNum++
			}
			break // stop circulation by key end
		}

		// is not key separator
		if _, ok := t.separators[key[right]]; !ok {
			right++
			continue
		}

		// is key separator
		segment := key[left : right+1]
		if child, ok := tmpRoot.Child[segment]; ok { // exist duplicate segment, find the child node
			tmpRoot.Size += size
			tmpRoot.KeyNum++
			tmpRoot = child
		} else {
			newNode := &Node{
				Segment: segment,
			}
			if tmpRoot.Child == nil {
				tmpRoot.Child = make(map[string]*Node)
			}
			tmpRoot.Child[segment] = newNode
			tmpRoot.Size += size
			tmpRoot.KeyNum++
			tmpRoot = newNode
			t.nodeNum++
		}
		left = right + 1
		right = left
	}
}

func (t *Tree) GetSize(keyPrefix string) int64 {
	size := int64(0)
	tmpRoot := t.root
	left := 0
	right := left
	for right < len(keyPrefix) {
		// key end
		if right == len(keyPrefix)-1 {
			if node, ok := tmpRoot.Child[keyPrefix[left:right+1]]; ok {
				size = node.Size
			}
			break
		}
		// is not key separator
		if _, ok := t.separators[keyPrefix[right]]; !ok {
			right++
			continue
		}
		// is key separator
		segment := keyPrefix[left : right+1]
		if child, ok := tmpRoot.Child[segment]; ok { // exist duplicate segment, find the child node
			tmpRoot = child
			left = right + 1
			right = left
		} else {
			right++ // even if not find segment split by separator, continue find longer segment
		}
	}
	return size
}

func (t *Tree) Expand(keyPrefix string) map[string]*Node {
	if keyPrefix == "" {
		return t.root.Child
	}

	tmpRoot := t.root
	left := 0
	right := left
	for right < len(keyPrefix) {
		// key end
		if right == len(keyPrefix)-1 {
			if node, ok := tmpRoot.Child[keyPrefix[left:right+1]]; ok {
				return node.Child
			}
			return nil
		}
		// is not key separator
		if _, ok := t.separators[keyPrefix[right]]; !ok {
			right++
			continue
		}
		// is key separator
		segment := keyPrefix[left : right+1]
		if child, ok := tmpRoot.Child[segment]; ok { // exist duplicate segment, find the child node
			tmpRoot = child
			left = right + 1
			right = left
		} else {
			right++ // even if not find segment split by separator, continue find longer segment
		}
	}
	return nil
}

func (t *Tree) GetKeyNum() int64 {
	return t.root.KeyNum
}

func (t *Tree) GetTotalSize() int64 {
	return t.root.Size
}

func (t *Tree) MergeSingleChildNode() {
	t.mergeSingleChildNode(t.root)
}

func (t *Tree) mergeSingleChildNode(node *Node) {
	if len(node.Child) == 0 {
		return
	}
	if len(node.Child) == 1 {
		for _, child := range node.Child {
			node.Segment += child.Segment
			node.Child = child.Child
			t.mergeSingleChildNode(node)
			break
		}
	} else {
		for _, child := range node.Child {
			t.mergeSingleChildNode(child)
		}
	}
}

// Print 打印字符串树
func (t *Tree) Print() {
	fmt.Printf("nodeNum:%v, keyNum:%v\n", t.nodeNum, t.root.KeyNum)
	t.root.print(0)
}

// print 打印节点
func (n *Node) print(level int) {
	fmt.Printf("%*s%s%*s%d\n", level*2, "", n.Segment, 2, "", n.Size)
	segments := make([]string, 0, len(n.Child))
	for segment := range n.Child {
		segments = append(segments, segment)
	}
	sort.Slice(segments, func(i, j int) bool {
		return segments[i] < segments[j]
	})
	for _, segment := range segments {
		n.Child[segment].print(level + 1)
	}
}

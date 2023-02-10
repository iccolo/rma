package analyzer

import (
	"fmt"
	"log"
	"sync"

	"github.com/iccolo/rma/internal/tree"
)

func NewKeyTypeTree(separators []byte) *KeyTypeTree {
	t := &KeyTypeTree{trees: [6]*tree.Tree{}}
	for i := 1; i <= 5; i++ {
		t.trees[i] = tree.New(KeyTypeToTypeStr[i], separators)
	}
	return t
}

type KeyTypeTree struct {
	trees [6]*tree.Tree
	rw    sync.RWMutex
}

func (k *KeyTypeTree) AddKey(info *KeyInfo) {
	k.rw.Lock()
	defer k.rw.Unlock()
	k.trees[info.KeyT].AddKey(info.Key, info.Size)
}

func (k *KeyTypeTree) GetSize(keyPrefix string, keyT KeyType) int64 {
	k.rw.RLock()
	defer k.rw.RUnlock()
	return k.trees[keyT].GetSize(keyPrefix)
}

func (k *KeyTypeTree) Expand(keyPrefix string, keyT KeyType) map[string]*tree.Node {
	k.rw.RLock()
	defer k.rw.RUnlock()
	return k.trees[keyT].Expand(keyPrefix)
}

func (k *KeyTypeTree) MergeSingleChildNode() {
	k.rw.Lock()
	defer k.rw.Unlock()
	for _, t := range k.trees {
		if t == nil {
			continue
		}
		t.MergeSingleChildNode()
	}
}

func (k *KeyTypeTree) Print() {
	fmt.Println("Summary:")
	for i, t := range k.trees {
		if t == nil {
			continue
		}
		fmt.Printf("Type:%s KeyNum:%d TotalSize:%d\n", KeyTypeToTypeStr[i], t.GetKeyNum(), t.GetTotalSize())
	}
	fmt.Println("Detail:")
	for _, t := range k.trees {
		if t == nil {
			continue
		}
		t.Print()
		fmt.Println()
	}
}

type KeyInfo struct {
	Key  string
	KeyT KeyType
	Size int64
}

func (a *Analyzer) analysisKey(keysChan chan []string, tree *KeyTypeTree, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		withTypeChan = make(chan []*KeyInfo, 100)
		withSizeChan = make(chan []*KeyInfo, 100)
	)
	wg.Add(3)
	go a.getKeyType(keysChan, withTypeChan, wg)
	go a.getKeySize(withTypeChan, withSizeChan, wg)
	go a.updateTree(withSizeChan, tree, wg)
}

func (a *Analyzer) updateTree(infoChan chan []*KeyInfo, tree *KeyTypeTree, wg *sync.WaitGroup) {
	defer wg.Done()

	var num int
	for infos := range infoChan {
		for _, info := range infos {
			tree.AddKey(info)
			num++
			if num%1000 == 0 {
				log.Printf("have analyze %v thousand keys\n", num/1000)
			}
		}
	}
	tree.MergeSingleChildNode()
	log.Println("analyze finish")
}

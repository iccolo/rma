package analyze

import (
	"container/heap"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/iccolo/rma/analyzer"
	"github.com/iccolo/rma/analyzer/tree"
)

type Handler interface {
	GetInstanceList() []*InstanceStatus
	StartAnalyze(ana *analyzer.Analyzer)
	GetKeyTypes(host string) ([]string, error)
	Expand(host, keyType, keyPrefix string, numLimit int64, sort SortVar) ([]*NodeInfo, error)
}

type handler struct {
	mu        sync.Mutex
	instances map[string]*instance
}

func NewHandler() Handler {
	h := &handler{
		instances: make(map[string]*instance),
	}
	return h
}

type instance struct {
	Host             string
	Tree             *analyzer.KeyTypeTree
	AnalyzeStartTime time.Time
	AnalyzeEndTime   time.Time
	IsFinish         bool
}

type InstanceStatus struct {
	Host             string `json:"host"`
	AnalyzeStartTime string `json:"analyze_start_time"`
	AnalyzeEndTime   string `json:"analyze_end_time"`
	IsFinish         bool   `json:"is_finish"`
}

func (h *handler) GetInstanceList() []*InstanceStatus {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.instances) == 0 {
		return nil
	}
	list := make([]*InstanceStatus, 0, len(h.instances))
	for _, instance := range h.instances {
		list = append(list, &InstanceStatus{
			Host:             instance.Host,
			AnalyzeStartTime: instance.AnalyzeStartTime.Format("2006-01-02 15:04:05"),
			AnalyzeEndTime:   instance.AnalyzeEndTime.Format("2006-01-02 15:04:05"),
			IsFinish:         instance.IsFinish,
		})
	}
	return list
}

func (h *handler) StartAnalyze(ana *analyzer.Analyzer) {
	//ana := &analyzer.Analyzer{
	//	Host:       a.Host,
	//	Port:       uint(a.Port),
	//	Password:   a.Password,
	//	Count:      uint(a.Count),
	//	Limit:      a.Limit,
	//	Match:      a.Match,
	//	Types:      strings.Split(a.Types, ","),
	//	Separators: []byte(a.Separators),
	//	Cluster:    a.Cluster,
	//	Pause:      time.Duration(a.Pause),
	//}
	t, wg := ana.AsyncRun()
	h.mu.Lock()
	h.instances[ana.Host] = &instance{
		Host:             ana.Host,
		Tree:             t,
		AnalyzeStartTime: time.Now(),
	}
	h.mu.Unlock()
	go func() {
		wg.Wait()
		h.mu.Lock()
		h.instances[ana.Host].AnalyzeEndTime = time.Now()
		h.instances[ana.Host].IsFinish = true
		h.mu.Unlock()
		log.Printf("finish analyze for host:%v", ana.Host)
	}()
	log.Printf("start analyze for host:%v", ana.Host)
}

func (h *handler) GetKeyTypes(host string) ([]string, error) {
	h.mu.Lock()
	instance, ok := h.instances[host]
	h.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("host:%v not exits", host)
	}
	return instance.Tree.GetKeyTypeStr(), nil
}

func (h *handler) Expand(host, keyType, keyPrefix string, numLimit int64, sortVar SortVar) ([]*NodeInfo, error) {
	h.mu.Lock()
	instance, ok := h.instances[host]
	h.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("host:%v not exits", host)
	}
	keyT, ok := analyzer.KeyTypeStrToType[keyType]
	if !ok {
		return nil, fmt.Errorf("req key type:%v not exist", keyType)
	}
	nodes := instance.Tree.Expand(keyPrefix, keyT)

	sortedNode := &SortedNode{
		Nodes:   make([]*tree.Node, 0, numLimit),
		SortVar: sortVar,
	}
	heap.Init(sortedNode)
	for _, node := range nodes {
		heap.Push(sortedNode, node)
		if sortedNode.Len() > int(numLimit) {
			heap.Pop(sortedNode)
		}
	}
	sort.Sort(sortedNode)
	layer := make([]*NodeInfo, 0, len(sortedNode.Nodes))
	for i := sortedNode.Len() - 1; i >= 0; i-- {
		node := sortedNode.Nodes[i]
		layer = append(layer, &NodeInfo{
			Segment:   node.Segment,
			KeyNum:    node.KeyNum,
			TotalSize: node.Size,
			ChildNum:  int32(len(node.Child)),
		})
	}
	return layer, nil
}

type NodeInfo struct {
	Segment   string `json:"segment"`
	KeyNum    int64  `json:"key_num"`
	TotalSize int64  `json:"total_size"`
	ChildNum  int32  `json:"child_num"`
}

type SortVar int32

const (
	SortVarTotalSize = 1
	SortVarKeyNum    = 2
	SortVarChildNum  = 3
)

type SortedNode struct {
	Nodes   []*tree.Node
	SortVar SortVar
}

func (e *SortedNode) Less(i, j int) bool {
	n := e.Nodes
	switch e.SortVar {
	case SortVarTotalSize:
		return n[i].Size < n[j].Size // Size 小优先Pop
	case SortVarKeyNum:
		return n[i].KeyNum < n[j].KeyNum // KeyNum 小优先
	case SortVarChildNum:
		return len(n[i].Child) < len(n[i].Child) // Child 数量少优先
	default:
		return n[i].Segment > n[j].Segment
	}
}

func (e *SortedNode) Swap(i, j int) {
	e.Nodes[i], e.Nodes[j] = e.Nodes[j], e.Nodes[i]
}

func (e *SortedNode) Len() int {
	return len(e.Nodes)
}

func (e *SortedNode) Push(x interface{}) {
	node := x.(*tree.Node)
	e.Nodes = append(e.Nodes, node)
}

func (e *SortedNode) Pop() interface{} {
	x := e.Nodes[e.Len()-1]
	e.Nodes = e.Nodes[:e.Len()-1]
	return x
}

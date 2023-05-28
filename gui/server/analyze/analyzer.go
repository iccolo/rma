package analyze

import (
	"container/heap"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/iccolo/rma/analyzer"
	"github.com/iccolo/rma/analyzer/tree"
)

type Handler interface {
	GetInstanceList() []*InstanceStatus
	StartAnalyze(ana *analyzer.Analyzer)
	GetKeyTypes(host string) ([]string, error)
	Expand(host, keyType, keyPrefix string, numLimit int64, sort SortVar) ([]*NodeInfo, error)
	GetKeyInfo(host, key string, limit int) (*RedisValue, error)
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
	Analyzer         *analyzer.Analyzer
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
	t, wg := ana.AsyncRun()
	h.mu.Lock()
	h.instances[ana.Host] = &instance{
		Host:             ana.Host,
		Tree:             t,
		AnalyzeStartTime: time.Now(),
		Analyzer:         ana,
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
		seg := node.Segment
		if len(node.Child) == 0 {
			seg = keyPrefix + seg
		}
		layer = append(layer, &NodeInfo{
			Segment:   seg,
			KeyNum:    node.KeyNum,
			TotalSize: node.Size,
			ChildNum:  int32(len(node.Child)),
		})
	}
	return layer, nil
}

func (h *handler) GetKeyInfo(host, key string, limit int) (*RedisValue, error) {
	h.mu.Lock()
	instance, ok := h.instances[host]
	h.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("host:%v not exits", host)
	}
	conn := instance.Analyzer.Dial()
	defer conn.Close()

	res := &RedisValue{
		Key: key,
	}
	// 查看 key 的类型
	keyType, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		return res, fmt.Errorf("failed to get key type: %s", err)
	}
	res.Type = keyType

	// 获取 key 的 TTL
	ttl, err := redis.Int64(conn.Do("TTL", key))
	if err != nil {
		return res, fmt.Errorf("failed to get key ttl: %s", err)
	}
	res.TTL = ttl

	// 根据 key 的类型获取对应的值
	switch keyType {
	case "string":
		// 获取 string 类型的值
		result, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return res, fmt.Errorf("failed to get string result: %s", err)
		}
		res.Value = result
	case "list":
		// 获取 list 类型的值
		results, err := redis.Strings(conn.Do("LRANGE", key, 0, limit-1))
		if err != nil {
			return res, fmt.Errorf("failed to get list results: %s", err)
		}
		res.Value = results
	case "hash":
		// 获取 hash 类型的值
		results, err := redis.Values(conn.Do("HSCAN", key, 0, "COUNT", limit))
		if err != nil {
			return res, fmt.Errorf("failed to get hash results: %s", err)
		}
		values, err := redis.StringMap(results[1], err)
		if err != nil {
			return res, fmt.Errorf("failed to get hash results: %s", err)
		}
		res.Value = values
	case "set":
		// 获取 set 类型的值
		results, err := redis.Values(conn.Do("SSCAN", key, 0, "COUNT", limit))
		if err != nil {
			return res, fmt.Errorf("failed to get set results: %s", err)
		}
		values, err := redis.Strings(results[1], err)
		if err != nil {
			return res, fmt.Errorf("failed to get set results: %s", err)
		}
		res.Value = values
	case "zset":
		// 获取 zset 类型的值
		results, err := redis.Values(conn.Do("ZRANGE", key, 0, limit-1, "WITHSCORES"))
		if err != nil {
			return res, fmt.Errorf("failed to get zset results: %s", err)
		}
		values := make([]ZSetItem, 0, len(results)/2)
		for i := 0; i < len(results); i += 2 {
			score, err := strconv.ParseFloat(string(results[i+1].([]byte)), 64)
			if err != nil {
				return res, fmt.Errorf("failed to parse zset score: %s", err)
			}
			values = append(values, ZSetItem{
				Member: string(results[i].([]byte)),
				Score:  score,
			})
		}
		res.Value = values
	default:
		return res, fmt.Errorf("unsupported redis key type '%s'", keyType)
	}

	return res, nil
}

type RedisValue struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	TTL   int64       `json:"ttl"`
	Value interface{} `json:"value"`
}

type ZSetItem struct {
	Member string  `json:"member"`
	Score  float64 `json:"score"`
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

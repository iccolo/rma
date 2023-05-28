package analyzer

import (
	"fmt"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type Analyzer struct {
	Host       string
	Port       uint `json:"port"`
	Password   string
	Count      uint   `json:"count"`
	Limit      uint64 `json:"limit"`
	Match      string
	Types      string
	Separators string
	Cluster    bool
	Pause      time.Duration `json:"pause"` // ms
}

func (a *Analyzer) Run() *KeyTypeTree {
	tree := NewKeyTypeTree([]byte(a.Separators))
	keysChan := make(chan []string, 10)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go a.scan(keysChan, wg)
	go a.analysisKey(keysChan, tree, wg)
	wg.Wait()
	return tree
}

func (a *Analyzer) AsyncRun() (*KeyTypeTree, *sync.WaitGroup) {
	tree := NewKeyTypeTree([]byte(a.Separators))
	keysChan := make(chan []string, 10)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go a.scan(keysChan, wg)
	go a.analysisKey(keysChan, tree, wg)
	return tree, wg
}

func (a *Analyzer) Dial() redigo.Conn {
	address := fmt.Sprintf("%s:%d", a.Host, a.Port)
	conn, err := redigo.Dial("tcp", address, redigo.DialPassword(a.Password))
	errorJudge("dial redis", err)
	return conn
}

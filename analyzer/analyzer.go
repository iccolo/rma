package analyzer

import (
	"fmt"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type Analyzer struct {
	Host       string
	Port       uint
	Password   string
	Count      uint
	Limit      uint64
	Match      string
	Types      []string
	Separators []byte
	Cluster    bool
	Pause      time.Duration // ms
}

func (a *Analyzer) Run() *KeyTypeTree {
	tree := NewKeyTypeTree(a.Separators)
	keysChan := make(chan []string, 10)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go a.scan(keysChan, wg)
	go a.analysisKey(keysChan, tree, wg)
	wg.Wait()
	return tree
}

func (a *Analyzer) dial() redigo.Conn {
	address := fmt.Sprintf("%s:%d", a.Host, a.Port)
	conn, err := redigo.Dial("tcp", address, redigo.DialPassword(a.Password))
	errorJudge("dial redis", err)
	return conn
}

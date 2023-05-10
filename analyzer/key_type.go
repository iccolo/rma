package analyzer

import (
	"strings"
	"sync"

	redigo "github.com/gomodule/redigo/redis"
)

type KeyType = int

const (
	KeyTypeString KeyType = 1
	KeyTypeList   KeyType = 2
	KeyTypeSet    KeyType = 3
	KeyTypeHash   KeyType = 4
	KeyTypeZset   KeyType = 5
)

func (a *Analyzer) getKeyType(keysChan chan []string, infoChan chan []*KeyInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	conn := a.dial()
	defer conn.Close()

	// set analyze key type, all types by default
	types := make(map[string]KeyType)
	for _, t := range strings.Split(a.Types, ",") {
		if kt, ok := KeyTypeStrToType[t]; ok {
			types[t] = kt
		}
	}
	if len(types) == 0 {
		types = KeyTypeStrToType
	}

	for keys := range keysChan {
		for _, key := range keys {
			err := conn.Send("TYPE", key)
			errorJudge("redis conn Send TYPE cmd", err)
		}
		err := conn.Flush()
		errorJudge("redis conn Flush", err)
		infos := make([]*KeyInfo, 0, len(keys))
		for _, key := range keys {
			result, err := redigo.String(conn.Receive())
			errorJudge("redis conn Receive", err)
			if keyT, ok := types[result]; ok {
				infos = append(infos, &KeyInfo{
					Key:  key,
					KeyT: keyT,
				})
			}
		}
		infoChan <- infos
	}
	close(infoChan)
}

var KeyTypeStrToType = map[string]KeyType{
	"string": KeyTypeString,
	"list":   KeyTypeList,
	"set":    KeyTypeSet,
	"hash":   KeyTypeHash,
	"zset":   KeyTypeZset,
}

var KeyTypeToTypeStr = map[KeyType]string{
	KeyTypeString: "string",
	KeyTypeList:   "list",
	KeyTypeSet:    "set",
	KeyTypeHash:   "hash",
	KeyTypeZset:   "zset",
}

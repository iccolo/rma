package analyzer

import (
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/iccolo/rma/analyzer/size"
)

func (a *Analyzer) getKeySize(inChan chan []*KeyInfo, outChan chan []*KeyInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	conn := a.Dial()
	defer conn.Close()

	if !a.Cluster {
		for infos := range inChan {
			for _, info := range infos {
				err := conn.Send("MEMORY USAGE", info.Key)
				errorJudge("redis conn Send MEMORY USAGE cmd", err)
			}
			err := conn.Flush()
			errorJudge("redis conn Flush", err)
			for _, info := range infos {
				result, err := redigo.Int64(conn.Receive())
				errorJudge("redis conn Receive", err)
				info.Size = result
			}
			outChan <- infos
			time.Sleep(a.Pause * time.Millisecond)
		}
		close(outChan)
		return
	}

	for infos := range inChan {
		for _, info := range infos {
			if f, ok := sendFunctions[info.KeyT]; ok {
				f(conn, info.Key)
			}
		}
		err := conn.Flush()
		errorJudge("redis conn Flush", err)
		for _, info := range infos {
			if f, ok := receiveFunctions[info.KeyT]; ok {
				members, length := f(conn)
				if ff, ok := sizeFunctions[info.KeyT]; ok {
					info.Size = int64(ff(info.Key, members, length))
				}
			}
		}
		outChan <- infos
		time.Sleep(a.Pause * time.Millisecond)
	}
	close(outChan)
	return
}

var sendFunctions = map[KeyType]func(conn redigo.Conn, key string){
	KeyTypeString: sendReadStringCmd,
	KeyTypeList:   sendReadListCmd,
	KeyTypeSet:    sendReadSetCmd,
	KeyTypeHash:   sendReadHashCmd,
	KeyTypeZset:   sendReadZsetCmd,
}

var receiveFunctions = map[KeyType]func(conn redigo.Conn) ([][]byte, int){
	KeyTypeString: receiveString,
	KeyTypeList:   receiveList,
	KeyTypeSet:    receiveSet,
	KeyTypeHash:   receiveHash,
	KeyTypeZset:   receiveZset,
}

var sizeFunctions = map[KeyType]func(string, [][]byte, int) int{
	KeyTypeString: size.String,
	KeyTypeList:   size.List,
	KeyTypeSet:    size.Set,
	KeyTypeHash:   size.Hash,
	KeyTypeZset:   size.Zset,
}

const sample = 5

func sendReadStringCmd(conn redigo.Conn, key string) {
	err := conn.Send("STRLEN", key)
	errorJudge("redis conn Send STRLEN cmd", err)
}

func sendReadListCmd(conn redigo.Conn, key string) {
	err := conn.Send("LLEN", key)
	errorJudge("redis conn Send LLEN cmd", err)
	err = conn.Send("LRANGE", key, 0, sample-1)
	errorJudge("redis conn Send STRLEN cmd", err)
}

func sendReadSetCmd(conn redigo.Conn, key string) {
	err := conn.Send("SCARD", key)
	errorJudge("redis conn Send SCARD cmd", err)
	err = conn.Send("SRANDMEMBER", key, sample)
	errorJudge("redis conn Send SRANDMEMBER cmd", err)
}

func sendReadHashCmd(conn redigo.Conn, key string) {
	err := conn.Send("HLEN", key)
	errorJudge("redis conn Send HLEN cmd", err)
	err = conn.Send("HSCAN", key, 0, "COUNT", sample)
	errorJudge("redis conn Send SSCAN cmd", err)
}

func sendReadZsetCmd(conn redigo.Conn, key string) {
	err := conn.Send("ZCARD", key)
	errorJudge("redis conn Send ZCARD cmd", err)
	err = conn.Send("ZRANGE", key, 0, sample-1)
	errorJudge("redis conn Send ZRANGE cmd", err)
}

func receiveString(conn redigo.Conn) ([][]byte, int) {
	length, err := redigo.Int(conn.Receive())
	errorJudge("redis Receive", err)
	return nil, length
}

func receiveList(conn redigo.Conn) ([][]byte, int) {
	length, err := redigo.Int(conn.Receive())
	errorJudge("redis Receive", err)
	members, err := redigo.ByteSlices(conn.Receive())
	errorJudge("redis Receive", err)
	return members, length
}

func receiveSet(conn redigo.Conn) ([][]byte, int) {
	length, err := redigo.Int(conn.Receive())
	errorJudge("redis Receive", err)
	members, err := redigo.ByteSlices(conn.Receive())
	errorJudge("redis Receive", err)
	return members, length
}

func receiveHash(conn redigo.Conn) ([][]byte, int) {
	length, err := redigo.Int(conn.Receive())
	errorJudge("redis Receive", err)
	results, err := redigo.Values(conn.Receive())
	errorJudge("redis Receive", err)
	memberValues, err := redigo.ByteSlices(results[1], nil)
	errorJudge("redis Receive", err)
	return memberValues, length
}

func receiveZset(conn redigo.Conn) ([][]byte, int) {
	length, err := redigo.Int(conn.Receive())
	errorJudge("redis Receive", err)
	members, err := redigo.ByteSlices(conn.Receive())
	errorJudge("redis Receive", err)
	return members, length
}

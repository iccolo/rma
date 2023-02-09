package analyzer

import (
	"log"
	"sync"

	redigo "github.com/gomodule/redigo/redis"
)

func (a *Analyzer) scan(keysChan chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn := a.dial()
	defer conn.Close()

	var (
		cursor int
		num    int
	)
	for {
		results, err := redigo.Values(conn.Do("SCAN", cursor, "MATCH", a.Match, "COUNT", a.Count))
		errorJudge("scan redis", err)
		cursor, _ = redigo.Int(results[0], nil)
		keys, _ := redigo.Strings(results[1], nil)
		num += len(keys)

		keysChan <- keys
		if cursor == 0 || uint64(num) >= a.Limit {
			break
		}
	}
	close(keysChan)
	log.Printf("scan finish, total %d keys\n", num)
}

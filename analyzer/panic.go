package analyzer

import "fmt"

func errorJudge(desc string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s:%v", desc, err))
	}
}

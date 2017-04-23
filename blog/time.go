package blog

import (
	"time"
	"fmt"
)

func RunTime(funcName string, startTime time.Time) {
	dis := time.Now().Sub(startTime).Seconds()
	fmt.Printf("Func Name:%s\nRun Time:%fs\n===\n", funcName, dis)
}
// @Author : Lik
// @Time   : 2021/7/14
package pprof

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

var dataList []string

func TestProfiling(t *testing.T) {
	go func() {
		for {
			log.Printf("len: %d", Add("go-programming-tour-book"))
			time.Sleep(time.Millisecond * 10)
		}
	}()

	_ = http.ListenAndServe(":8080", nil)
}

func Add(str string) int {
	data := []byte(str)
	dataList = append(dataList, string(data))
	return len(dataList)
}

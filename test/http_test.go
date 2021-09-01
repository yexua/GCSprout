// @Author : Lik
// @Time   : 2021/1/27
package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("pong")
	})

	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("恢复系统")
			}
		}()
		panic("我出错了")
	})

	http.ListenAndServe(":8000", nil)

}

// Package channel
// @Author : Lik
// @Time   : 2021/9/1
package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("发送数据")
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println("a: ", a)
		}
	}()
	compare()
	time.Sleep(time.Second * 2)
	close(ch)
	fmt.Println("ok")

}

func compare() {
	sm1 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	sm2 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	if reflect.DeepEqual(sm1, sm2) {
		fmt.Println("sm1 == sm2")
	}
}

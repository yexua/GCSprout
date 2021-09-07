// Package sync
// @Author : Lik
// @Time   : 2021/9/6
package main

import (
	"fmt"
	"sync"
)

func main() {
	once()

}

func printing() {
	ch := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			ch <- 1
			if i%2 == 1 {
				fmt.Println("A:", i)
			}
		}
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			<-ch
			if i%2 == 0 {
				fmt.Println("B:", i)
			}
		}
	}()

	for {
	}
}

func once() {
	var o sync.Once
	o.Do(func() {
		var a = 1 / 0
		fmt.Println(a)

	})
}

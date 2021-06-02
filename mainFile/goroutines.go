package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	resMap = make(map[int]int, 10)

	// 对map读写操作进行加锁
	lock sync.Mutex
)

func factorial(n int) {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}

	resMap[n] = res

}

func main() {
	for i := 1; i <= 10; i++ {
		go factorial(i)
	}

	for k, v := range resMap {
		fmt.Println(k, v)
	}
}

func print() {
	for i := 0; i < 10; i++ {
		fmt.Println("print() hello,Word", i)
		time.Sleep(time.Second)
	}
}

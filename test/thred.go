package main

import (
	"errors"
	"fmt"
	"time"
)

func mainss() {
	randomNums := randomNum(1, 100000, 100000)
	a := heapSort(randomNums, 100)
	fmt.Println(a) //打印这最大的100个

	ch := make(chan int, 2)

	// 发送1个数据关闭channel
	ch <- 1
	ch <- 0
	close(ch)
	print("close channel\n")

	// 不停读数据直到channel没有有效数据
	for {
		select {
		case v, ok := <-ch:
			print("v: ", v, ", ok:", ok, "\n")
			if !ok {
				print("channel is close\n")
				return
			}
		default:
			print("nothing\n")
		}
	}
}

func f3() func() int {
	a := 10086
	return func() int {
		return a
	}
}

// 只有generator进行对outCh进行写操作，返回声明
// <-chan int，可以防止其他协程乱用此通道，造成隐藏bug
func generator(n int) <-chan int {
	outCh := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			outCh <- i
		}
	}()
	return outCh
}

// consumer只读inCh的数据，声明为<-chan int
// 可以防止它向inCh写数据
func consumer(inCh <-chan int) {
	for x := range inCh {
		fmt.Println(x)
	}
}

func doWithTimeOut(timeout time.Duration) (int, error) {
	select {
	case ret := <-do():
		return ret, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

func do() <-chan int {
	outCh := make(chan int)
	go func() {
		// do work
	}()
	return outCh
}

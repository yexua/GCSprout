// @Author : Lik
// @Time   : 2021/1/27
package error

import (
	"fmt"
	"testing"
	"time"
)

// 跨携程失效
func TestPanic(t *testing.T) {

	defer println("In main")
	go func() {
		defer println("in goroutine")
		panic("this is a panic")
	}()

	time.Sleep(time.Second)
}

func TestRecover(t *testing.T) {
	defer println("in main")
	if err := recover(); err != nil {
		println(err)
	}

	panic("unknown err")
}

// 嵌套崩溃
func TestMorePanic(t *testing.T) {

	defer println("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		panic("panic again")
	}()

	panic("panic once")
}

func TestPanicContinue(t *testing.T) {
	defer func() {
		fmt.Println("恢复panic")
		if err := recover(); err != nil {
			fmt.Println("恢复完毕")
		}

	}()
	println(1)
	println(2)
	println(3)
	println(4)
	panic("出错了")
	println(5)
}

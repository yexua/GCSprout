// @Author : Lik
// @Time   : 2021/1/26
package sync

import (
	"runtime"
	"sync"
)

func TestSchedDemo() {
	wg := new(sync.WaitGroup)
	wg.Add(4)
	go func() {
		defer wg.Done()
		for i := 0; i < 6; i++ {
			println(i)
			if i == 3 {
				runtime.Gosched()
			}
		}
	}()
	go func() {
		defer wg.Done()
		println("Hello, World - 1")
	}()
	go func() {
		defer wg.Done()
		println("Hello, World - 2")
	}()
	go func() {
		defer wg.Done()
		println("Hello, World - 3")
	}()

	wg.Wait()
}

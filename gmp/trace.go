// Package gmp
// @Author : Lik
// @Time   : 2021/8/31
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

func main() {
	// 创建一个trace文件
	file, err := os.Create("trace.out")
	if err != nil {

	}
	defer file.Close()

	// 启动trace
	err = trace.Start(file)
	if err != nil {
		panic(err)
	}

	// 要调试的业务
	fmt.Println(runtime.GOMAXPROCS(0))
	fmt.Println("Hello trace")

	// 停止trace
	trace.Stop()

}

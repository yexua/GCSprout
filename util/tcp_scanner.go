// @Author : Lik
// @Time   : 2021/2/2
package util

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func scannerSingle(addr string, start, end uint32) {

	for i := start; i < end; i++ {
		address := fmt.Sprintf("%s:%d", addr, i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("%s 关闭了\n", address)
			continue
		}
		conn.Close()
		fmt.Printf("%s 打开了！！！\n", address)
	}
}

func scanner并发(addr string, start, end uint32) {
	var wg sync.WaitGroup
	for i := start; i < end; i++ {
		wg.Add(1)
		go func(j uint32) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", addr, i)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf("%s 关闭了\n", address)
				return
			}
			conn.Close()
			fmt.Printf("%s 打开了！！！\n", address)
		}(i)
	}
	wg.Wait()
}

func worker(addr string, ports chan int, results chan int) {
	//defer wg.Done()
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- -1
			continue
		}
		conn.Close()
		results <- p
	}

}

func Scanner并发池(addr string, start, end uint32) {
	s := time.Now()
	ports := make(chan int, 100)
	// main 需要使用，只有自己，不需要缓冲
	results := make(chan int)
	var openPorts []int
	var closePorts []int
	//var wg sync.WaitGroup
	for i := 0; i < cap(ports); i++ {
		go worker(addr, ports, results)
	}

	go func() {
		for i := start; i < end; i++ {
			fmt.Printf("开始扫描端口: %d\n", i)
			ports <- int(i)
		}
	}()

	len := int(end - start)
	for i := 0; i < len; i++ {
		port := <-results
		if port != -1 {
			openPorts = append(openPorts, port)
		} else {
			closePorts = append(closePorts, port)
		}
	}
	//wg.Wait()
	fmt.Println("执行完毕")
	close(ports)
	close(results)

	sort.Ints(openPorts)

	for i := range openPorts {
		fmt.Printf("%d 打开了\n", openPorts[i])
	}

	elapsed := time.Since(s) / 1e9
	fmt.Printf("speed %d seconds\n", elapsed)

}

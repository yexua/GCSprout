package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	n := 0
	replay := &n
	Multiply(2, 3, replay)
	fmt.Println("Multiply", *replay)

	x := min(1, 3, 2, 0)
	fmt.Printf("The minimum is: %d\n", x)
	slice := []int{7, 9, 3, 5, 1}
	x = min(slice...)
	fmt.Printf("The minimum in the slice is: %d\n", x)
	x = min_array([]int{-1, 2, 3, 4})
	fmt.Printf("The minimum is: %d\n", x)


	a()

	fmt.Println(f())
}

func Multiply(a, b int, reply *int) {
	*reply = a * b
}

//变长参数
func min(s ...int) int {
	if len(s) == 0 {
		return 0
	}
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

func min_array(s []int) int {
	if len(s) == 0 {
		return 0
	}
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

// defer 关键字
func a() {
	i := 0
	defer fmt.Println(i)
	i++
	return
}


// 延迟表达式必须是函数调用
func func1(s string) (n int, err error) {
	defer func() {
		log.Printf("func1(%q) = %d, %v", s, n, err)
	}()
	return 7, io.EOF
}


func f() (ret int) {
	defer func() {
		ret++
	}()
	return 1
}

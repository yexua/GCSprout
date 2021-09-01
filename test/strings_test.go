// @Author : Lik
// @Time   : 2021/1/22
package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestString(t *testing.T) {

	var buf strings.Builder
	len, err := buf.WriteString("我爱你")
	if err != nil {

	}
	fmt.Println(len)
	fmt.Println(buf.String())
}

func TestChan(t *testing.T) {
	c := make(chan int, 1)
	c <- 2.0
}

func TestPointer(t *testing.T) {
	var f float64
	fmt.Println(reflect.TypeOf(&f))
	fmt.Println(reflect.TypeOf(unsafe.Pointer(&f)))
	fmt.Println(reflect.TypeOf((*int64)(unsafe.Pointer(&f))))

}

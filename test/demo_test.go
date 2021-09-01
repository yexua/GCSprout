// @Author : Lik
// @Time   : 2021/1/28
package main

import (
	"GCSprout/util"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestArray(t *testing.T) {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	for i := 0; i < 3; i++ {
		fmt.Println(&arr[i])
	}
	for _, v := range arr {
		fmt.Println(&v)
		newArr = append(newArr, &v)
	}
	for _, v := range newArr {
		fmt.Println(*v)
	}

}

func TestTcp(t *testing.T) {
	util.Scanner并发池("8.131.92.162", 1, 1024)
}

func TestDate(t *testing.T) {
	now := time.Now().Local()
	date := now.AddDate(0, -1, 0)

	fmt.Println(now)
	fmt.Println(date)
}

func TestCompared(t *testing.T) {
	sm1 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	sm2 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	//if sm1 == sm2 {
	//	fmt.Println("sm1 == sm2")
	//}
	if reflect.DeepEqual(sm1, sm2) {

	}
}

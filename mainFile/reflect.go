package main

import (
	"fmt"
	"reflect"
)

type val []interface{}

type Student struct {
	Name string
	Age  int
}

func reflectInterface(obj interface{}) {

	// 先获取到reflect.Type
	rtype := reflect.TypeOf(obj)
	fmt.Println("type:", rtype)

	// 获取到 reflect.value
	rVal := reflect.ValueOf(obj)
	fmt.Printf("rVal=%v rVal type is %T\n", rVal, rVal)

}

func reflectStruct(obj interface{}) {
	rtype := reflect.TypeOf(obj)
	fmt.Println("type:", rtype)

	// 获取到 reflect.value
	rVal := reflect.ValueOf(obj)
	iv := rVal.Interface()
	fmt.Printf("rVal=%v rVal type is %T\n", iv, iv)

	// 将interface{} 通过断言转成需要的类型
	if stu, ok := iv.(Student); ok {
		fmt.Printf("name: %v Age: %v", stu.Name, stu.Age)
	}
}

func main() {

	var num int = 100
	reflectInterface(num)

	stu := Student{
		Name: "Tom",
		Age:  20,
	}

	reflectStruct(stu)

}

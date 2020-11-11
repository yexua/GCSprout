package main

import (
	"fmt"
	"math/rand"
	"time"
)

//构建小顶堆
func min(root int, end int, a []int) {
	for {
		var child = 2*root + 1
		//判断是否存在child节点
		if child > end {
			break
		}
		//判断右child是否存在，如果存在则和另外一个同级节点进行比较
		if child+1 <= end && a[child] > a[child+1] {
			child++
		}
		if a[root] > a[child] {
			a[root] = a[child]
			a[child] = a[root]
			root = child
		} else {
			break
		}
	}
}

//再次构建堆
func generateHeap(arr []int, end int) {
	for start := end / 2; start >= 0; start-- {
		min(start, end, arr)
	}
}

//在a数组中找出num个最大值
func heapSort(a []int, num int) []int {
	l := len(a) - 1
	generateHeap(a[:num], num-1)
	fmt.Println("已生成")
	for i := num; i <= l; i++ {
		if a[0] < a[i] {
			a[0] = a[i]
			a[i] = a[0]
			generateHeap(a[:num], num-1)
		}
	}
	fmt.Println(a[:num])
	return a
}

//随机数组
//start是开始的数 end是结束的数 其中start>0 end>start
func randomNum(start int, end int, count int) []int {
	//给定范围
	if end < start || (end-start) < count {
		return nil
	}
	//存放随机后的结果
	nums := make([]int, 0)
	t := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(nums) < count {
		num := t.Intn((end - start))
		exist := false
		for _, v := range nums {
			if v == num { //说明生成了一样的
				exist = true
				break
			}
		}
		if exist == false {
			nums = append(nums, num)
		}
	}
	return nums
}

func f1() (r *int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return &t
}

func mains() {

	fmt.Println(*f1())

	randomNums := randomNum(1, 100, 10)
	fmt.Println("生成", randomNums)
	a := heapSort(randomNums, 1)
	fmt.Println(a)
}

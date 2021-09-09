package main

import (
	"GCSprout/zap"
	"bufio"
	"container/heap"
	"fmt"
	"github.com/gin-gonic/gin"
	uberZap "go.uber.org/zap"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	findMid([]int{10, 11, 21, 19, 21, 17, 21, 18, 15})
	findMid([]int{1, 2, 2, 2, 3, 4, 3, 4, 4, 5})

	fmt.Println(Median([]int{1, 2, 3, 4}))

	//reverseWords()
	//findMin2("12345264", 4)
}

func serve() {
	gin.DisableConsoleColor()
	r := gin.New()
	logger, _ := uberZap.NewProduction()
	defer logger.Sync()
	r.Use(zap.GinZapMiddleware(logger, "2006-01-02 15:04:05", false))
	//r.Use(zap.RecoveryWithZap(logger, true))
	r.GET("/ping", func(c *gin.Context) {

		time.Sleep(time.Second * 2)
		logger.Info("业务执行成功")

		c.JSON(http.StatusOK, gin.H{
			"res": "Hello World",
		})
	})

	r.GET("/dlv", TestDlvHandler)

	r.GET("/panic", func(c *gin.Context) {
		panic("This is a panic!")
	})

	r.Run(":9090")
}

func TestDlvHandler(c *gin.Context) {
	fmt.Println("接受到请求")
	a := 1
	b := a + 1
	d := b + 2
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(d)

	c.JSON(http.StatusOK, gin.H{
		"res": "Successes!",
	})
}

func compareVersion(version1 string, version2 string) int {

	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	s1 := len(v1)
	s2 := len(v2)
	if s1 > s2 {
		v2 = append(v2, make([]string, s1-s2)...)
	} else if s1 < s2 {
		v1 = append(v1, make([]string, s2-s1)...)
	}

	for i := 0; i < s1; i++ {
		index1, _ := strconv.Atoi(v1[i])
		index2, _ := strconv.Atoi(v2[i])
		if index1 > index2 {
			return 1
		} else if index1 < index2 {
			return -1
		}
	}
	return 0
}

func compareVersion2(version1 string, version2 string) int {
	l1, l2 := len(version1), len(version2)
	i, j := 0, 0
	for i < l1 || j < l2 {
		x := 0
		for ; i < l1 && version1[i] != '.'; i++ {
			x = x*10 + int(version1[i]-'0')
		}
		i++
		y := 0
		for ; j < l2 && version2[j] != '.'; j++ {
			y = y*10 + int(version2[j]-'0')
		}
		j++
		if x > y {
			return 1
		}

		if x < y {
			return -1
		}
	}

	return 0
}

// 最长不重复子字符串长度
func lengthOfLongestSubstring(s string) int {
	max := 0

	size := len(s)
	for i := 0; i < size; i++ {
		m := make(map[uint8]byte)
		for j := i; j < size; j++ {
			if _, ok := m[s[j]]; ok {
				break
			} else {
				m[s[j]] = 0
			}
		}
		if len(m) > max {
			max = len(m)
		}
	}
	return max
}

func lengthOfLongestSubstring2(s string) int {
	//m := make(map[byte]byte)
	//
	//max, start := 0, 0
	//for end := 0; end < len(s); end++ {
	//
	//	if _, ok := m[s[end]]; !ok {
	//		m[s[end]] = 0
	//		if m[s[end]] + 1 > start{
	//			start = m[s[end]]
	//		}
	//	}
	//
	//	if map.containsKey(ch) {
	//		start = Math.max(map.get(ch)+1, start)
	//	}
	//	max = Math.max(max, end-start+1)
	//	map.put(ch, end)
	//}
	//return max
	return 0
}

func testSearch() {
	nums := []int{-1, 0, 3, 5, 9, 12}
	res := search(nums, 2)
	fmt.Println("找到了:", res)
}

func search(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		//m := (left + right) / 2
		m := left + (right-left)/2
		if nums[m] > target {
			right = m - 1
		} else if nums[m] < target {
			left = m + 1
		} else {
			return m
		}
	}
	return -1
}

func searchRange(nums []int, target int) []int {
	res := []int{-1, -1}

	l := binarySearch(nums, target)
	r := binarySearch(nums, target+1) - 1
	if l == len(nums) || nums[l] != target {
		return res
	} else {
		res[0] = l
		res[1] = l
		if r > l {
			res[1] = r
		}
	}

	return res
}

// binarySearch 二分查找
func binarySearch(nums []int, target int) int {
	left, right := 0, len(nums)
	for left < right {
		m := left + (right-left)/2
		if nums[m] >= target {
			right = m
		} else {
			left = m + 1
		}
	}
	return left
}

func balancedStringSplit(s string) int {
	res, count := 0, 0
	if len(s) <= 0 {
		return res
	}
	for _, v := range s {
		if v == 'L' {
			count += 1
		} else {
			count -= 1
		}
		if count == 0 {
			res += 1
		}
	}
	return res
}

// https://leetcode-cn.com/problems/ipo/
func findMaximizedCapital(k int, w int, profits []int, capital []int) int {
	n := len(profits)
	arr := make([][2]int, n)
	for i := 0; i < n; i++ {
		arr[i][0] = capital[i]
		arr[i][1] = profits[i]
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i][0] < arr[j][0]
	})

	h := &hp{
		bool: true,
	}
	for cur := 0; k > 0; k-- {
		for cur < n && arr[cur][0] <= w {
			heap.Push(h, arr[cur][1])
			cur++
		}
		if h.Len() == 0 {
			break
		}
		w += heap.Pop(h).(int)
	}
	return w
}

func testFindMaximizedCapital() {
	findMaximizedCapital(2, 0, []int{1, 2, 3}, []int{0, 1, 1})
}

type hp struct {
	sort.IntSlice
	bool
}

// Less 默认小顶堆，bool为true则为大顶堆
func (h hp) Less(i, j int) bool {

	if h.bool {
		return h.IntSlice[i] > h.IntSlice[j]
	}
	return h.IntSlice[i] < h.IntSlice[j]

}

func (h *hp) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}

func (h *hp) Pop() interface{} {
	a := h.IntSlice
	v := a[len(a)-1]
	h.IntSlice = a[:len(a)-1]
	return v
}

func (h *hp) Top() interface{} {
	return h.IntSlice[h.Len()-1]
}

func reverseWords() {

	start, end := 0, 0
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')

	fmt.Scanln(&start)
	fmt.Scanln(&end)

	res := ""
	s = strings.Trim(s, "\n")
	strs := strings.Split(s, " ")
	var word []string
	for i := 0; i < len(strs); i++ {
		if strs[i] != "" {
			word = append(word, strs[i])
		}
	}
	for i := 0; i < len(word); i++ {
		if i == start && i < end {
			word[i], word[end] = word[end], word[i]
			start++
			end--
		}
		res = res + word[i] + " "
	}
	fmt.Print(res)
}

// 查找一组数据中出现次数最多的数，可以是多个，组成新数组并求出中位数
func findMid(arr []int) {
	m := make(map[int]int)
	max := 1
	var nArr []int
	for _, v := range arr {
		if _, ok := m[v]; ok {
			m[v]++
			if m[v] > max {
				max = m[v]
				nArr = nArr[:0]
			}
		} else {
			m[v] = 1
		}
		if m[v] == max {
			nArr = append(nArr, v)
		}
	}

	size := len(nArr)
	hpSize := size/2 + 1

	hp := &hp{}

	for i := 0; i < hpSize; i++ {
		hp.Push(nArr[i])
	}

	for i := hpSize; i < size; i++ {
		if nArr[i] > hp.Top().(int) {
			hp.Pop()
			hp.Push(nArr[i])
		}
	}
	var result float64
	if size%2 == 0 {
		result = float64(hp.Pop().(int)+hp.Top().(int)) / 2
	} else {
		result = float64(hp.Top().(int))
	}
	fmt.Println(result)
}

func Median(arr []int) float64 {
	hp := &hp{bool: false}
	size := len(arr)
	hpSize := size/2 + 1
	for i := 0; i < hpSize; i++ {
		hp.Push(arr[i])
	}

	for i := hpSize; i < size; i++ {
		if arr[i] > hp.Top().(int) {
			hp.Pop()
			hp.Push(arr[i])
		}
	}
	var result float64
	if size%2 == 0 {
		a := hp.Pop().(int)
		b := hp.Top().(int)
		result = float64(a+b) / 2.0
	} else {
		result = float64(hp.Top().(int))
	}
	return result
}

// 输入的第一行为一个字符串，字符串由0-9字符组成，记录正整数num1，num1的长度小于32
// 输入的第二行为需要移除的数字的个数，小于num1的长度
// 输入：2615371
//       4
// 输出：131
func findMin() {
	s := ""
	n := 0
	fmt.Scanln(&s)
	fmt.Scanln(&n)

	size := len(s)
	arr := make([]int, size)
	for i, v := range s {
		d, _ := strconv.Atoi(string(v))
		arr[i] = d
	}

	res := make([]int, size-n)
	for i, j := n, 0; j < size-n; j++ {
		res[j] = arr[i]
		i++
	}

	min := -1
	last := -1

	for i := 0; i < size-n; i++ {
		flag := true
		for j := n - 1 + i; j > last; j-- {
			if res[i] >= arr[j] {
				res[i] = arr[j]
				min = j
				flag = false
			}
		}
		last = min
		if flag {
			break
		}
	}
	for i := range res {
		fmt.Print(i)
	}
}

func findMin2(s string, n int) {
	stack := make([]int32, 0, len(s))
	for _, c := range s {
		var i int
		for i = len(stack) - 1; i >= 0 && stack[i] > c && n > 0; i-- {
			n--
		}
		stack = stack[:i+1]
		stack = append(stack, c)
	}
	ans := strings.TrimLeft(string(stack), "0")
	if ans == "" {
		ans = "0"
	}
	fmt.Println(ans)
}

type LinkNode struct {
	val  int
	next *LinkNode
}

// 判断一个链表是否有环
func test(head *LinkNode) {
	var fast *LinkNode
	var slow *LinkNode

	fast, slow = head, head

	for slow != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
		if slow == fast {
			fmt.Println("存在环")
		}
	}
	fmt.Println("不存在环")
}

func gro() {

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			fmt.Println("hello")
		}()
	}
	wg.Wait()
	fmt.Println("world")
}

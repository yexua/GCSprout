// @Author : Lik
// @Time   : 2020/10/21
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var points = [][]int{{3, 3}, {5, -1}, {-2, 4}, {4, 6}}
	kClosest(points, 2)
	for i := range points {
		fmt.Println(points[i])
	}
}
func kClosest(points [][]int, K int) (ans [][]int) {
	rand.Shuffle(len(points), func(i, j int) {
		points[i], points[j] = points[j], points[i]
	})

	var quickSelect func(left, right int)
	quickSelect = func(left, right int) {
		if left == right {
			return
		}
		pivot := points[right]
		lessCount := left
		for i := left; i < right; i++ {
			if less(points[i], pivot) {
				points[i], points[lessCount] = points[lessCount], points[i]
				lessCount++
			}
		}
		points[right], points[lessCount] = points[lessCount], points[right]
		if lessCount+1 == K {
			return
		} else if lessCount+1 < K {
			quickSelect(lessCount+1, right)
		} else {
			quickSelect(left, lessCount-1)
		}
	}
	quickSelect(0, len(points)-1)
	return points[:K]
}

func less(p, q []int) bool {
	return p[0]*p[0]+p[1]*p[1] < q[0]*q[0]+q[1]*q[1]
}

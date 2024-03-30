package main

import (
	"fmt"
	"sort"
)

//func minimumSum(nums []int) int {
//	n := len(nums)
//	const inf = 1 << 30
//	right := make([]int, n+1)
//	right[n] = inf
//	for i := n - 1; i >= 0; i-- {
//		right[i] = min(right[i+1], nums[i])
//	}
//	ans, left := inf, inf
//	for i, x := range nums {
//		if left < x && right[i] < x {
//			ans = min(ans, left+x+right[i])
//		}
//		left = min(left, x)
//	}
//	if ans == inf {
//		return -1
//	}
//	return ans
//}

func main() {
	//nums =[8,6,1,5,3]
	nums := []int{8, 6, 1, 5, 3}
	sort.Ints(nums)
	fmt.Println(nums)
}
func minimumSum(nums []int) int {
	const inf = 1 >> 30
	n := len(nums)
	right := make([]int, n+1)
	right[n] = inf
	for i := n - 1; i >= 0; i-- {
		right[i] = min(right[i+1], nums[i])
	}
	ans, left := inf, inf
	for i, num := range nums {
		if left < num && num > right[i] {
			ans = min(ans, left+num+right[i])
		}
		left = min(left, num)
	}
	if ans == inf {
		return -1

	}
	return ans
}

package main

func minimumSum(nums []int) int {
	var min int
	min = nums[0]
	index := 0
	for i, num := range nums {
		if num < min {
			min = num
			index = i
		}
	}
	if index < 2 {
		s := make([]int, 3, 5)
		index1 := len(nums) - 1
		rmin := nums[index1]
		for i := index1; i > index; i-- {
			if rmin > nums[i] {
				rmin = nums[i]
				index1 = i
			}
		}
		for i := index1 - 1; i > index; i-- {
			if nums[i] > nums[index1] {
				s = append(s, nums[i])
			}
		}
		var mmin int
		mmin = s[0]
		for _, i := range s {
			if mmin > i {
				mmin = i
			}
		}
		return min + mmin + rmin
	}
	if index > len(nums)-3 {
		s := make([]int, 3, 5)
		index1 := 0
		lmin := nums[0]
		for i := 0; i < index; i++ {
			if lmin > nums[i] {
				lmin = nums[i]
				index1 = i
			}
		}
		for i := index1 + 1; i < index; i++ {
			if nums[i] > nums[index1] {
				s = append(s, nums[i])
			}
		}
		var mmin int
		mmin = s[0]
		for _, i := range s {
			if mmin > i {
				mmin = i
			}
		}
		return min + mmin + lmin
	}
	if index >= 2 && index <= len(nums)-3 {
		s := make([]int, 3, 5)
		index1 := len(nums) - 3
		rmin := nums[index1]
		for i := index1; i > index; i-- {
			if rmin > nums[i] {
				rmin = nums[i]
				index1 = i
			}
		}
		if index == index1 {
			if nums[index+1] > nums[index+2] {
				return min + nums[index+1] + nums[index+2]
			}
		}
		for i := index1 - 1; i > index; i-- {
			if nums[i] > nums[index1] {
				s = append(s, nums[i])
			}
		}
		var mmin int
		mmin = s[0]
		for _, i := range s {
			if mmin > i {
				mmin = i
			}
		}
		if index == index1 {
			var val3, val4 int
			if nums[index+1] > nums[index+2] {
				val3 = min + nums[index+1] + nums[index+2]
			} else {
				return 1000000
			}
			if nums[index-1] > nums[index-2] {
				val4 = min + nums[index-1] + nums[index-2]
			} else {
				return 1000000
			}
			return max(val3, val4)
		} else {
			return -1
		}
		val1 := min + mmin + rmin
		s1 := make([]int, 3, 5)
		index2 := 0
		lmin := nums[0+1]
		for i := 0; i < index; i++ {
			if lmin > nums[i] {
				lmin = nums[i]
				index2 = i
			}
		}
		for i := index2 + 1; i < index; i++ {
			if nums[i] > nums[index2] {
				s1 = append(s, nums[i])
			}
		}
		var wmin int
		wmin = s1[0]
		for _, i := range s {
			if wmin > i {
				wmin = i
			}
		}
		val2 := min + wmin + lmin
		return max(val2, val1)
	}
	return -1
}
func main() {
	//nums =[8,6,1,5,3]
	nums := []int{8, 6, 1, 5, 3}
	minimumSum(nums)
}

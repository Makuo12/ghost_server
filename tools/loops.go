package tools

import "fmt"


func NumInAscendingOrder(nums []float64) (numsOrdered map[int]float64) {
	var currentSmallNum float64 = nums[0]
	var currentSmallNumIndex int = 0
	loopLength := len(nums)-1
	for i := 0; i < loopLength; i++ {
		for idx, u :=range nums{
			if u < currentSmallNum {
				currentSmallNum = u
				currentSmallNumIndex = idx
			}
		}
		nums = append(nums[:currentSmallNumIndex], nums[currentSmallNumIndex+1:]...)
		
		numsOrdered[i] = currentSmallNum
		fmt.Println(nums)
		if len(nums) <=  1{
			numsOrdered[i] = nums[0]
			nums = nums[0+1:]
		}else{
			fmt.Println(len(nums), i)
			currentSmallNum = nums[len(nums)-1]
			currentSmallNumIndex = len(nums)-1
		}
		
	}
	fmt.Println(nums)
	return numsOrdered
}

// Takes a list and add none to the list if the list is empty
func HandleDBList(a []string) []string {
	if len(a) == 0 {
		return []string{"none"}
	}
	return a
}


func ConcatSlices(slice1, slice2 []int) []int {
	result := make([]int, len(slice1)+len(slice2))
	copy(result, slice1)
	copy(result[len(slice1):], slice2)
	return result
}
package utils

func MinMax(nums []int) (int, int) {

    if len(nums) == 0 {
        return 0, 0
    }

    min := nums[0]
    max := nums[0]

    for i:= 1; i<len(nums); i++ {
        if nums[i] < min {
            min = nums[i]
        }
        if nums[i] > max {
            max = nums[i]
        }
    }
    
    return min, max
}
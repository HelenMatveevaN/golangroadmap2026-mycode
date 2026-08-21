//2 pointers
package main

func twoSum(numbers []int, target int) []int {
	left := 0
	right := len(numbers) - 1

	for left < right {
		currentSum := numbers[left] + numbers[right]

		if currentSum == target {
			return []int{left+1, right+1}
		} else if currentSum < target {
			left++
		} else if currentSum > target {
			right--
		}
	}

	return[]int{}
}
//2 pointers
package main

func sortedSquares(nums []int) []int {
	n := len(nums)
	result := make([]int, n)

	left := 0
	right := n - 1

	for i := n-1; i >= 0; i -- {
		leftSquare := nums[left] * nums[left]
		rightSguqre := nums[right] * nums[right]

		if leftSquare > rightSguqre {
			result[i] = leftSquare
			left++
		} else {
			result[i] = rightSguqre
			right--
		}
	}

	return result
}
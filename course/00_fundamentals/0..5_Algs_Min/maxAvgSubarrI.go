//Скользящее окно Sliding Window
package main

func findMaxAverage(nums []int, k int) float64 {
	//curr win
	currentSum := 0
	for i := 0; i < k; i++ {
		currentSum += nums[i]
	}

	maxSum := currentSum

	for i := k; i < len(nums); i++ {
		currentSum += nums[i] - nums[i-k]

	if currentSum > maxSum {
		maxSum = currentSum
		}
	}

	// Среднее значение — это максимальная сумма, деленная на k.
	return float64(maxSum) / float64(k)
}
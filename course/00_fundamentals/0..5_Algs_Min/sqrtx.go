
//Bynary Search -1
package main

func mySqrt(x int) int {
	//для 0 и 1
	if x < 2 { 
		return x
	}

	left := 1
	right := x / 2
	ans := 0

	for left <= right {
		mid := left + (right - left)/2

		if mid <= x / mid {
			ans = mid //возможный ответ
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return ans
}
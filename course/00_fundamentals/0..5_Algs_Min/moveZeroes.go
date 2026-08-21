//2 pointers
package main

func moveZeroes(nums []int) {
	lastNonZeroesFoundAt := 0 //медл.указ-ль

	for cur := 0; cur < len(nums); cur++ { //быстр.указ-ль, перебирает весь массив
		if nums[cur] != 0 {
			if cur != lastNonZeroesFoundAt {
				nums[cur], nums[lastNonZeroesFoundAt] = nums[lastNonZeroesFoundAt], nums[cur]
			}
			lastNonZeroesFoundAt++
		}
	}
}
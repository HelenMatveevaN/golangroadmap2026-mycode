//Массивы, Хэш-таблицы и База

//Given an integer array nums, return true if any value 
//appears at least twice in the array, and return false if every element is distinct.

package main

func containsDuplicate(nums []int) bool {
    seen := make(map[int]struct{})

    for _, num := range nums {
    	if _, exists := seen[num]; exists {
    		return true //дубликат найден
    	}
    	//не нашли
    	seen[num] = struct{}{}
    }

    return false
}
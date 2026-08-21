//Скользящее окно Sliding Window
package main

func containsNearbyDuplicate(nums []int, k int) bool {
    // Хэш-таблица для хранения элементов текущего окна
    window := make(map[int]bool)

    for i, num := range nums {
    	// 1. Проверяем, есть ли такое число в текущем окне
    	if window[num] {
    		return true
    	}

    	// 2. Добавляем текущее число в окно
    	window[num] = true

    	if len(window) > k {
    		oldestNum := nums[i-k]
    		delete(window, oldestNum)
    	}
    }

    return false
}
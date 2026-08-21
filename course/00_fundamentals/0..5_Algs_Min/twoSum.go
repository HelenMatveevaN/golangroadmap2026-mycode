//Массивы, Хэш-таблицы и База

/*
Вам дан массив целых чисел nums и целое число target. 
Верните индексы этих двух чисел так, чтобы их сумма равнялась
target .
*/

package main

//nums = [2,7,11,15], target = 9
func twoSum(nums []int, target int) []int {
	// Хэш-таблица: ключ — само число, значение — его индекс в массиве
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num //недостающее число до суммы

		if idx, found := numMap[complement]; found {
			//нашли
			return []int{idx, i}
		}

		//не нашли, запоминаем число и индекс
		numMap[num] = i
	}

	return []int{}
}
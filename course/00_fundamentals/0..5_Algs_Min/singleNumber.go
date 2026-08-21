func singleNumber(nums []int) int {
    result := 0

    for _, num := range nums {
    	result ^= num //применяем XOR к каждому числу
    }

    return result
}
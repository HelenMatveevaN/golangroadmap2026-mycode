//Скользящее окно Sliding Window
package main


func countGoodSubstrings(s string) int {
    count := 0
    n := len(s)

    if n < 3 {
    	return 0
    }

    for i :=0; i <= n-3; i++ {
    	b1 := s[i]
    	b2 := s[i+1]
    	b3 := s[i+2]

    	//проверка на уникальность
    	if b1 != b2 && b2 != b3 && b3 != b1 {
    		count++
    	}
    }

    return count
}
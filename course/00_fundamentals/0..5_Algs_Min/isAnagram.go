func isAnagram(s string, t string) bool {

//Given two strings s and t, 
//return true if t is an anagram of s, and false otherwise.	

	if len(s) != len(t) {
		return false
	}

	counts := make(map[rune]int)

	for _, r := range s { //идем посимвольно
		counts[r]++
	}
	for _, r := range t {
		counts[r]--
		if counts[r] < 0 {
			return false
		}
	}

	return true
}
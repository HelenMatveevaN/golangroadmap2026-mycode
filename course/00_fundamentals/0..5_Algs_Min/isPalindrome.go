//2 pointers
//s consists only of printable ASCII characters.

package main

func isPalindrome(s string) bool {
	left := 0
	right := len(runes)-1

	for left < right {
		//пропускаем не-буквы-цифры
		if !isAlphanumeric(s[left]) {
			left++
			continue
		}
		if !isAlphanumeric(s[right]) {
			right--
			continue
		}

		if toLower(s[left]) != toLower(s[right]) {
			return false
		}

		//сдвигаем оба указателя к центру
		left++
		right--
	}

	return true
}

// Проверка ASCII символа без использования тяжелого пакета unicode
func isAlphanumeric(b byte) byte {
	return (b >= 'a' && b <= 'z') ||
		   (b >= 'A' && b <= 'Z') ||
		   (b >= '0' && b <= '9')
}

// Приведение ASCII символа к нижнему регистру математически
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32 // В таблице ASCII строчные буквы сдвинуты на 32 относительно заглавных
	}
	return b
}
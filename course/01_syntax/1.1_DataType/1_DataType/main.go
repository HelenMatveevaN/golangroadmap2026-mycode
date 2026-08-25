package main

import (
	"fmt"
	"math"
	"unsafe"
	"time"
)

// 1. Конвертация из int64 в int8 с проверкой переполнения сверху и снизу
func int64ToInt8(val int64) (int8, error) {
	if val > math.MaxInt8 || val < math.MinInt8 {
		return 0, fmt.Errorf("overflow: value %d bounds out of int8 [%d, %d]", val, math.MinInt8, math.MaxInt8)
	}
	return int8(val), nil
}	

// 2. Конвертация из int в uint с проверкой потери отрицательного знака
func intToUint(val int) (uint, error) {
	if val < 0 {
		return 0, fmt.Errorf("overflow/underflow: negative value %d cannot be uint", val)
	}
	return uint(val), nil
}

// 3. Конвертация из float64 в int64 (потеря дробной части не считается overflow, 
// но выход за границы максимального целого числа — да)
func float64ToInt64(val float64) (int64, error) {
	if val > math.MaxInt64 || val < math.MinInt64 {
		return 0, fmt.Errorf("overflow: float %f out of int64 range", val)
	}
	return int64(val), nil
}

func main() {
	fmt.Println("--- Тест 1: Системное (тихое) переполнение в Go ---")
	var max8 int8 = math.MaxInt8 //127
	fmt.Printf("Максимум int8: %d\n", max8)

	// Прибавляем 1 к максимальному числу. Рантайм молчит, значение превращается в минимальное.
	fmt.Printf("Максимум int8 + 1 (тихое переполнение): %d\n\n", max8+1)

	fmt.Println("--- Тест 2: Безопасное преобразование int64 -> int8 ---")
	validInt64	:= int64(42)
	overflowInt64 := int64(150)

	if res, err := int64ToInt8(validInt64); err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Успешно: %d (тип %T)\n", res, res)
	}

	if _, err := int64ToInt8(overflowInt64); err != nil {
		fmt.Println("Поймали ошибку:", err)
	}

	fmt.Println("\n--- Тест 3: Преобразование со сменой знака int -> uint ---")
	negativeInt := -5
	if _, err := intToUint(negativeInt); err != nil {
		fmt.Println("Поймали ошибку:", err)
		// Что было бы при неявном / неосторожном касте:
		fmt.Printf("Что вернул бы обычный uint(-5): %d\n", uint(negativeInt))		
	}

	fmt.Println("\n--- Тест 4: Преобразование float64 -> int64 ---")
	hugeFloat := 9.5e19 // Очень большое число, заведомо больше MaxInt64
	if _, err := float64ToInt64(hugeFloat); err != nil {
		fmt.Println("Поймали ошибку:", err)
	}

	fmt.Println("\n--- Тест 5: Эксперименты ---")
	var v1 int
	fmt.Printf("Размер int: %d байт\n", unsafe.Sizeof(v1))

	var t time.Time
	fmt.Println("Значение времени:", t)
	if t.IsZero() {
		fmt.Println("IsZero вернул true")
	}

	var a int = 10
	var b int32 = 10
	//var c int32 = a
	_ = a
	_ = b	
	//_ = c
}
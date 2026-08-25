package main

import "fmt"

func ProcessDynamicData(input any) {
	switch v := input.(type) {
	case int:
		fmt.Printf("Это целое число * 2: %d\n", v*2)
	case string:
		fmt.Printf("Длина строки: %d\n", len(v))
	case bool:
		fmt.Println("Логический тип: ", v)	
	default:
		fmt.Printf( "Неизвестный тип: %T\n", v)	
	}
}

func main() {
	ProcessDynamicData(5)
	ProcessDynamicData("три")
	ProcessDynamicData(false)
	ProcessDynamicData(5.)
}
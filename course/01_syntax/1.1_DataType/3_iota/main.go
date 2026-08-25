package main

import "fmt"

type DeliveryStatus uint8

/*
Интерфейс Stringer находится в пакете fmt и требует всего один метод:

type Stringer interface {
    String() string
}
*/

const (
	StatusUnknown DeliveryStatus = iota 	//0
	StatusInWarehouse 					//1 (посылка на складе).
	StatusInTransit 					//2 (посылка в пути)
	StatusInTruck 						//3 (курьер везет посылку)
	StatusDelivered 					//4 (доставлено)
)

// Реализация интерфейса Stringer
func (s DeliveryStatus) String() string {
	switch s {
	case StatusUnknown:
		return "Неизвестный статус"
	case StatusInWarehouse:
		return "На складе"
	case StatusInTransit:
		return "В пути"
	case StatusInTruck:
		return "Курьер везет"
	case StatusDelivered:
		return "Доставлено"
	default:
		return fmt.Sprintf("Неизвестный код (99)")
	}
}

func main() {
	// Создаем переменную без явного указания значения
	var initialStatus DeliveryStatus
	// Проверка Zero Value
	fmt.Printf("Zero value статус: %v\n", initialStatus)

	initialStatus = 2
	fmt.Printf("New статус: %v\n", initialStatus)

	var BrokenStatus DeliveryStatus = DeliveryStatus(10)
	fmt.Printf("Сломанный статус: %v\n", BrokenStatus)

}
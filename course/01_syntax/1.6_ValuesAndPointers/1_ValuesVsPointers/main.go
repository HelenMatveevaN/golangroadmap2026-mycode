package main

//Сравнить производительность value vs pointer receiver для большой структуры.

// BigStruct имитирует "тяжелый" объект (например, конфиг или документ).
// Наш массив из 1280 элементов int64 займет ровно 10 240 байт (~10 КБ).
type BigStruct struct {
	data [1280]int64
}

// ValueReceiverMethod принимает копию всех 10 КБ данных в стек.
//go:noinline // Запрещаем компилятору оптимизировать (встраивать) метод
func (b BigStruct) ValueReceiverMethod() int64 {
	return b.data[0]
}

// PointerReceiverMethod принимает только 8 байт адреса в памяти.
//go:noinline
func (b *BigStruct) PointerReceiverMethod() int64 {
	return b.data[0]
}
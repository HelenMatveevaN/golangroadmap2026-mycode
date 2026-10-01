package main

import (
	"fmt"
)

/*
Конвейер из 3х стадий: 

Генератор (Stage 1): Генерирует числа от 1 до 5 и 
отправляет их в первый канал.

Удвоитель (Stage 2): Читает из первого канала, 
умножает каждое число на 2 
и отправляет во второй канал.

Квадратор (Stage 3): Читает из второго канала, 
возводит число в квадрат 
и отправляет в финальный канал, 
который вычитывается в main.
*/

func generator(nums ...int) <-chan int {
	ch := make(chan int)

	go func(nn []int, ch chan int){
		for _, num := range nn {
			ch <- num
		}
		close(ch)
	}(nums, ch)
	
	return ch
}

func multiplier(in <-chan int) <-chan int {
	out := make(chan int)

	go func(){
		for num := range in {
			out <- num * 2
		}
		close(out) //здесь in закрыт
	}()

	return out
}

func squarer(in <-chan int) <-chan int {
	out := make(chan int)

	go func(){
		for num := range in {
			out <- num * num
		}
		close(out) //здесь in закрыт
	}()

	return out
}

func main() {
	results := squarer(multiplier(generator(1,2,3,4,5)))

	for res := range results {
		fmt.Println(res)
	}
}
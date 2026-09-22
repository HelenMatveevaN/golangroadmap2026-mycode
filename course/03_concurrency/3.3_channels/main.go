package main

/*Практика: пинг-понг между двумя горутинами через канал.*/

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int //сч-к ударов
}

func main() {
	table := make(chan *Ball)

	go player("Пинг 🏓", table)
	go player("Понг 🔔", table)

	table <- new(Ball)
	time.Sleep(2 * time.Second)
	<-table

	fmt.Println("Игра окончена!")
}

func player(name string, table chan *Ball) {
	for {
		ball, ok := <-table
		if !ok {
			return
		}

		ball.hits++
		fmt.Printf("%s: удар №%d\n", name, ball.hits)

		time.Sleep(200 * time.Millisecond)

		table <- ball
	}
}
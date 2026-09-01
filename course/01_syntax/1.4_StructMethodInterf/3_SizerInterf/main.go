package main

import (
	"fmt"
	"math"
)

//Реализация интерфейса: 
//Написать интерфейс Sizer с методом Area() float64. 
//Реализовать его для структур Rectangle (прямоугольник) и Circle (круг).

type Sizer interface {
	Area() float64
}

type Rectangle struct {
	Width float64
	Heigh float64
}

type Circle struct {
	Radius float64
}

// метод Area для Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Heigh
}

// метод Area для Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

//Универсальная функция
func PrintArea(s Sizer) {
	fmt.Printf("Площадь фигуры: %.2f\n", s.Area())
}

func main() {
	rect  := Rectangle{Width: 10, Heigh: 5}
	circl := Circle{Radius: 3}

	fmt.Print("Прямоугольник — ")
	PrintArea(rect)

	fmt.Print("Круг — ")
	PrintArea(circl)
}
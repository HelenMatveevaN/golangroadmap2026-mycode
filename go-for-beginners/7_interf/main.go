package main

import (
    "fmt"
    "math"
)

// Shape — интерфейс для геометрических фигур
type Shape interface {
    Area() float64          // Метод для вычисления площади
    Perimeter() float64     // Метод для вычисления периметра
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rect struct {
    W, H float64
}

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

//Интерфейсы в Go — неявные. Тип реализует интерфейс, просто имея нужные методы. 
//Никаких implements.

//Задание: добавь тип Triangle и реализуй для него интерфейс Shape.
// Triangle — структура, представляющая треугольник по трем его сторонам.
type Triangle struct {
    SideA float64
    SideB float64
    SideC float64
}

func (t Triangle) Perimeter() float64 { return t.SideA + t.SideB + t.SideC }
func (t Triangle) Area() float64      { 
    p := t.Perimeter() //полупериметр
    //формула Герона: S = sqrt(p * (p - a) * (p - b) * (p - c))
    return math.Sqrt(p * (p - t.SideA) * (p - t.SideB) * (p - t.SideC))
}

func printShape(s Shape) {
    fmt.Printf("Площадь (%T): %.2f, Периметр: %.2f\n", s, s.Area(), s.Perimeter())
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rect{W: 4, H: 6},
        Triangle{SideA: 2, SideB: 3, SideC: 4},
    }
    
    for _, s := range shapes {
        printShape(s)
    }
}
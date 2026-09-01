package main

import (
    "fmt"
    "errors"
)

type Amount struct {
    Units   int64  //целая часть: рубли, доллары
    Cents   int64  //дробная часть: копейки, центы
}

type Money struct {
    Amount      //embedded struct
    Currency    string
}

//метод Add(other Money) (Money, error) с проверкой совпадения валют.
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("валюты не совпадают")
    }

    newUnits := m.Units + other.Units
    newCents := m.Cents + other.Cents

    if newCents >= 100 {
        newUnits += newCents / 100
        newCents = newCents % 100
    }

    return Money{
        Amount:     Amount{Units: newUnits, Cents: newCents},
        Currency:   m.Currency,
    }, nil
}

//Реализовать метод Mul(factor int64) Money для умножения суммы на число
func (m Money) Mul(factor int64) Money {
    newUnits := factor * m.Units
    newCents := factor * m.Cents

    if newCents >= 100 {
        newUnits += newCents / 100
        newCents = newCents % 100
    }

    return Money{
        Amount:     Amount{Units: newUnits, Cents: newCents},
        Currency:   m.Currency,
    }
}

//Реализовать интерфейс Stringer (метод String() string) 
//для красивого вывода (например, 100.50 RUB).
func (m Money) String() string {
    return fmt.Sprintf("%d.%02d %s", m.Units, m.Cents, m.Currency)
}

func main() {
    item1 := Money{
        Amount:     Amount{Units: 100, Cents: 50},
        Currency:   "RUB",
    }

    item2 := Money{
        Amount:     Amount{Units: 50, Cents: 70},
        Currency:   "RUB",
    }

    fmt.Println("Товар 1: ", item1)
    fmt.Println("Товар 2: ", item2)

    fmt.Println("\n--- Тест 1: Сложение ---")
    total, err := item1.Add(item2)
    if err != nil {
        fmt.Println("Ошибка при сложении: ", err)
    } else {
        fmt.Println("Итого сумма: ", total)
    }

    fmt.Println("\n--- Тест 2: Умножение ---")
    expensiveItems := item1.Mul(3)
    fmt.Println("Товар 1 в тройном экземпляре:", expensiveItems)
}
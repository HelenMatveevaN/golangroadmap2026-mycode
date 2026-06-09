package main

import (
    "errors"
    "fmt"
    "math"
    "sync"
)

type Point struct {
    X, Y float64
}

// Метод с получателем-значением
func (p Point) String() string {
    return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

// Метод с указателем — может изменять структуру
func (p *Point) Scale(factor float64) {
    p.X *= factor
    p.Y *= factor
}

func (p Point) Distance(other Point) float64 {
    dx := p.X - other.X
    dy := p.Y - other.Y
    return math.Sqrt(dx*dx + dy*dy)
}

type Rectangle struct {
    TopLeft     Point
    BottomRight Point
}

func (r Rectangle) Area() float64 {
    w := math.Abs(r.BottomRight.X - r.TopLeft.X)
    h := math.Abs(r.BottomRight.Y - r.TopLeft.Y)
    return w * h
}

//=============================================================
//Задание: создай структуру BankAccount 
//с методами Deposit, Withdraw (с проверкой баланса) и Balance.

type BankAccount struct {
    mu              sync.Mutex
    AccountNumber   string
    balance         float64
}

//Deposit - пополняет баланс счета на указанную сумму
func (ba *BankAccount) Deposit(amount float64) error {
    ba.mu.Lock()
    defer ba.mu.Unlock()

    //блок валидации
    if amount <=0 {
        return errors.New("сумма пополнения должна быть больше 0")
    }

    ba.balance += amount
    return nil
}

//Withdraw - списывает указанную сумму с баланса счета
func (ba *BankAccount) Withdraw(amount float64) error {
    ba.mu.Lock()
    defer ba.mu.Unlock()

    //блок валидации
    if amount <=0 {
        return errors.New("сумма пополнения должна быть больше 0")
    }

    if ba.balance - amount < 0 {
        return errors.New("недостаточно средств на счету для списания")
    }

    ba.balance -= amount
    return nil
}

//Balance - проверка баланса
func (ba *BankAccount) Balance() float64 {
    ba.mu.Lock()
    defer ba.mu.Unlock()

    return ba.balance
}

//=============================================================

func main() {
    p := Point{3, 4}
    fmt.Println(p)
    p.Scale(2)
    fmt.Println(p)

    origin := Point{0, 0}
    fmt.Printf("Расстояние: %.2f\n", p.Distance(origin))

    rect := Rectangle{Point{0, 0}, Point{5, 3}}
    fmt.Printf("Площадь: %.1f\n", rect.Area())

    //Задание: создай структуру BankAccount с методами Deposit, Withdraw (с проверкой баланса) и Balance.
    fmt.Println("=========== Задание: создай структуру BankAccount с методами... =======")

    account := &BankAccount{
        AccountNumber: "40817810000001234567",
        balance:       300.15,
    }

    //Пополнение
    err := account.Deposit(1000)
    if err == nil {
        fmt.Printf("Баланс после пополнения: %.2f\n", account.Balance())
    } else {
        fmt.Println("Ошибка при пополнении:", err)
    }

    //Списание
    err = account.Withdraw(-400)
    if err == nil {
        fmt.Printf("Баланс после списания: %.2f\n", account.Balance())
    } else {
        fmt.Println("Ошибка при списании:", err)
    }
}
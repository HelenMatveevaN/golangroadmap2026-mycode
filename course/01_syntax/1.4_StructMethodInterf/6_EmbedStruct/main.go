package main

import (
    "fmt"
)

//Использовать embedded struct для композиции
type User struct {
    Name string
}

type Admin struct {
    User
    Role string
}

func (u User) Greet() {
    fmt.Printf("Привет, меня зовут %s!\n", u.Name)
}

func main() {
    myAdmin := Admin{
        User: User{Name: "Алиса"},
        Role: "Супер-админ",
    }

    myAdmin.Greet()

    fmt.Println("Имя админа:", myAdmin.User)
}
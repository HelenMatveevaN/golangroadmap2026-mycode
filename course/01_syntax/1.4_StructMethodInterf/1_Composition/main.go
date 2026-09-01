package main

//Композиция: Создать структуру User и встроить её в структуру Admin. 
//Посмотреть, как Admin автоматически получает доступ к полям User

import "fmt"

// 1. Базовая структура User
type User struct {
	ID 		int
	Name 	string
}

// 2. Структура Admin, в которую мы ВСТРАИВАЕМ структуру User
type Admin struct {
	User //встраивание
	Role	string
}

func main() {
	myAdmin := Admin{
		User: User{
			ID:		1,
			Name:	"Алексей",
		},
		Role: "Superadmin",
	}

	fmt.Println("Имя админа:", myAdmin.Name) // Выведет: Алексей (а не myAdmin.User.Name)
	fmt.Println("ID админа:", myAdmin.ID)     // Выведет: 1

	fmt.Println("Роль админа:", myAdmin.Role) // Выведет: Superadmin

	fmt.Println("Полный путь до имени:", myAdmin.User.Name) //полный путь

}
/*
Теги JSON: json:"name" — имя поля, 
omitempty — пропустить если пустое, 
"-" — никогда не сериализовывать.

Задание: 
добавь эндпоинт GET /users/{id} с возвратом пользователя по ID 
и ответом 404, если не найден.
*/

package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "time"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email,omitempty"`
    Password string `json:"-"`
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "Алексей", Email: "alex@example.com"},
        {ID: 2, Name: "Мария"},
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "неверный JSON", http.StatusBadRequest)
        return
    }
    user.ID = 42
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    //имитация бд из шаблона
    users := []User{
        {ID: 1, Name: "Алексей", Email: "alex@example.com"},
        {ID: 2, Name: "Мария"},
    }

    //вытащить id пути
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "invalid ID format"})
        return
    }

    //ищем пользователя в срезе
    for _, user := range users {
        if user.ID == id {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(user)
            return
        }
    }

    //цикл кончился, пользователя нет
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/users", usersHandler)
    mux.HandleFunc("/users/create", createUserHandler)
    //NEW
    mux.HandleFunc("/users/", getUserByIDHandler) // слэш в конце ловит users/1, ...

    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    fmt.Println("Сервер запущен на :8080")
    if err := server.ListenAndServe(); err != nil {
        fmt.Println("Ошибка:", err)
    }
}
package main

import (
    "fmt"
    "runtime"
    "time"
    "unsafe"
)

// RemoveWithOrder удаляет элемент и сохраняет порядок.
// Сложность: O(n) — рантайму приходится двигать все элементы справа влево.
func RemoveWithOrder(slice []int, index int) []int {
    // Валидация: если индекс за границами, возвращаем слайс без изменений
    if index < 0 || index >= len(slice) {
        return slice
    }

    // Склеиваем кусок ДО индекса (slice[:index]) 
    // и кусок ПОСЛЕ индекса (slice[index+1:])
    return append(slice[:index], slice[index+1:]...)
}

// RemoveFast удаляет элемент за O(1), меняя его местами с последним.
func RemoveFast(slice []int, index int) []int {
    if index < 0 || index >= len(slice) {
        return slice
    }

    slice[index] = slice[len(slice)-1]
    slice[len(slice)-1] = 0

    return slice[:len(slice)-1]
}

// RemoveByValue удаляет ВСЕ элементы, равные target, без выделения новой памяти.
// Сложность: O(n) по времени, O(1) по памяти (ноль аллокаций).
func RemoveByValue(slice []int, target int) []int {
    // Создаем новый Slice Header, который смотрит на НАЧАЛО того же массива
    res := slice[:0]

    for _, v := range slice {
        if v != target {
            res = append(res, v) //перезаписывает память
        }
    }
    return res
}

type User struct{Name string}

// RemovePointerWithOrder удаляет указатель и зануляет старую ячейку для GC.
func RemovePointerWithOrder(slice []*User, index int) []*User {
    if index < 0 || index >= len(slice) {
        return slice
    }

    // Сдвигаем элементы влево, затирая удаляемый элемент
    copy(slice[index:], slice[index+1:])

    // ВАЖНО: Зануляем самый последний элемент, который теперь дублируется
    slice[len(slice)-1] = nil

    return slice[:len(slice)-1]
}

func main() {
    slice := []int{10, 20, 30, 40, 50, 60}

    slice = RemoveWithOrder(slice, 1)
    fmt.Println("RemoveWithOrder:", slice)

    slice = RemoveFast(slice,2)
    fmt.Println("RemoveFast:", slice)


    fmt.Println("\n--- ТЕСТ RemoveByValue ---")
    // 1. Создаем исходный слайс с дубликатами тройки
    source := []int{1, 3, 2, 3, 4, 3, 5}
    
    // Запоминаем адрес начала базового массива до операции
    originalAddr := unsafe.Pointer(&source[0])
    originalCap := cap(source)

    fmt.Printf("До удаления:   %v (len: %d, cap: %d, addr: %p)\n", 
        source, len(source), originalCap, originalAddr)

    // 2. Удаляем все тройки
    result := RemoveByValue(source, 3)

    // Запоминаем адрес результата
    resultAddr := unsafe.Pointer(&result[0])

    fmt.Printf("После удаления: %v (len: %d, cap: %d, addr: %p)\n", 
        result, len(result), cap(result), resultAddr)

    // 3. Проверки (Assertions)
    fmt.Println("\n--- Итоги проверки ---")
    if originalAddr == resultAddr {
        fmt.Println("✅ Память сохранена! Новый слайс использует ТОТ ЖЕ базовый массив.")
    } else {
        fmt.Println("❌ Ошибка! Произошла аллокация новой памяти.")
    }

    if cap(result) == originalCap {
        fmt.Println("✅ Емкость (capacity) осталась неизменной.")
    }

    // 4. Побочный эффект Sharing Memory (Важно увидеть!)
    // Поскольку мы работали in-place, первые элементы исходного слайса перезаписались
    fmt.Printf("⚠️ Что стало с исходным слайсом source: %v\n", source)


    fmt.Println("\n--- ТЕСТ RemovePointerWithOrder ---")
    // 1. Создаем три объекта в куче
    u1 := &User{Name: "Alice"}
    u2 := &User{Name: "Bob"}
    u3 := &User{Name: "Charlie"}

    // 2. Вешаем "маячки" (Finalizers). 
    // Как только GC уничтожит объект, мы увидим принт в консоли.
    runtime.SetFinalizer(u1, func(u *User) { fmt.Printf("🔥 GC удалил: %s\n", u.Name) })
    runtime.SetFinalizer(u2, func(u *User) { fmt.Printf("🔥 GC удалил: %s\n", u.Name) })
    runtime.SetFinalizer(u3, func(u *User) { fmt.Printf("🔥 GC удалил: %s\n", u.Name) })

    // 3. Формируем слайс указателей
    users := []*User{u1, u2, u3}

    // 4. Стираем внешние локальные ссылки, чтобы объекты держались ТОЛЬКО внутри слайса
    u1, u2, u3 = nil, nil, nil    

    fmt.Println("--- Удаляем Bob (индекс 1) ---")
    users = RemovePointerWithOrder(users, 1)

    // 5. Форсированно запускаем сборщик мусора
    runtime.GC()

    // Даем горутине GC крошечную паузу, чтобы успеть напечатать логи
    time.Sleep(100 * time.Millisecond)

    fmt.Printf("Осталось в слайсе (len: %d): ", len(users))
    for _, u := range users {
        fmt.Printf("%s ", u.Name)
    }
    fmt.Println("\n--- Конец теста ---")    

}
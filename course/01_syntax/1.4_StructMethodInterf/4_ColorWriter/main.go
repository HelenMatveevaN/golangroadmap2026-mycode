package main

import (
    "fmt"
    "os"
)

/*
//интерфейс под капотом
type Writer interface {
    Write(p []byte) (n int, err error)
}
*/

// ColorWriter будет оборачивать вывод и красить текст
type ColorWriter struct {
    target *os.File // Сюда мы запишем os.Stdout (консоль)
}

// Write реализует контракт интерфейса io.Writer
func (cw *ColorWriter) Write(p []byte) (n int, err error) {
    // Включаем зеленый цвет
    _, err = cw.target.WriteString("\033[32m")
    if err != nil {
        return 0, err
    }

    n, err = cw.target.Write(p)
    if err != nil {
        return n, err
    }

    // Выключаем цвет
    _, err = cw.target.WriteString("\033[0m")
    if err != nil {
        return n, err
    }

    //возвращаем количество записанных байт и nil (ошибок нет)
    return n, nil
}

func main() {
    myWriter := &ColorWriter{target: os.Stdout}

    //Функция fmt.Fprintln не знает про наш ColorWriter
    // Так как у нашего *ColorWriter есть метод Write, компилятор разрешает нам его передать

    fmt.Fprintln(myWriter, "Привет! Этот текст будет напечатан зеленым цветом!")
    fmt.Fprintln(myWriter, "И эта строчка тоже автоматически покрасится!")
}
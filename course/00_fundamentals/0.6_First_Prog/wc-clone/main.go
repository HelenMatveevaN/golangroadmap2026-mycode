package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите имя файла")
		return
	}
	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Ошибка: %g\n", err)
		return
	}

	//подсчет cnt /n
	lineCount := bytes.Count(data, []byte("\n"))

	fmt.Printf("%d %s\n", lineCount, filename)
}

/*
➜  wc-clone git:(main) ✗ wc -l code_copy.txt
      25 code_copy.txt

➜  wc-clone git:(main) ✗ go run main.go code_copy.txt
25 code_copy.txt

*/
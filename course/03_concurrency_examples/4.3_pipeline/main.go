package main

import (
	"fmt"
	"strings"
)

func generate(words ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, w := range words {
			out <- w
		}
	}()
	return out
}

func toUpper(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- strings.ToUpper(s)
		}
	}()
	return out
}

func addExclamation(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- s + "!"
		}
	}()
	return out
}

func main() {
	// Конвейер: generate → toUpper → addExclamation
	words := generate("hello", "world", "golang")
	upper := toUpper(words)
	excited := addExclamation(upper)

	for result := range excited {
		fmt.Println(result) // HELLO!, WORLD!, GOLANG!
	}
}
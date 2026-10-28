package main

import (
	"fmt"
	"context"
	"flag"
	"os"
	"os/signal"
	"time"
)

//check timeout + Cntrl+C
//~go run main.go

//go run main.go -timeout=10s

func main() {
	//УКАЗАТЕЛЬ на time.Duration (*time.Duration)
	var timeoutFlag = flag.Duration("timeout", 5*time.Second, "timeout for the operstion")

	//set vars from os.Args
	flag.Parse() // *timeoutFlag

	//родительский контекст
	baseCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt) //os.Interrupt (Ctrl+C)
	defer stop()

	//дочерний контекст
	ctx, cancel := context.WithTimeout(baseCtx, *timeoutFlag)
	defer cancel()

	fmt.Printf("Программа запущена с таймаутом %v. Нажмите Ctrl+C для выхода...\n", *timeoutFlag)

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				fmt.Println("Причина: Время истекло (Timeout).")
			case context.Canceled:
				fmt.Println("Причина: Программа прервана пользователем (Ctrl+C).")
			}
			return

		case t := <-ticker.C:
			// Имитация работы
			fmt.Printf("[%s] Какая-то полезная работа выполняется...\n", t.Format("15:04:05"))
		}
	}
}

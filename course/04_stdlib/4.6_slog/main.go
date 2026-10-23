package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
	//настройка логгера
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger) //сделала дефолтным

	start := time.Now()

	//имитация работы
	time.Sleep(50 * time.Millisecond) 	

	requestID := "req-xyz123"
	userID := 42
	latencyMS := time.Since(start).Milliseconds()

	/*slog.Info("Запрос успешно обработан",
		"request_id", requestID,
		"user_id", userID,
		"latency_ms", latencyMS,
	)*/

	slog.Info("Запрос успешно обработан",
		slog.String("request_id", requestID),
		slog.Int("user_id", userID),
		slog.Int64("latency_ms", latencyMS),
	)
}
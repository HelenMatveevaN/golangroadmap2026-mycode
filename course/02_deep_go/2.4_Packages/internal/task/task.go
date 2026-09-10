package task

/*
Здесь лежат только доменная модель и ошибки, которые к ней относятся
*/

import "errors"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Глобальные ошибки валидации и бизнес-логики
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
)
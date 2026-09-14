package calc

import (
	"errors"
	"strconv"
	"strings"
	
	// Импортируем дженерик-контейнеры
	"my-module/golangroadmap2026/course/02_deep_go/containers"
)

// Operation — интерфейс на стороне потребителя
type Operation interface {
	Execute(a, b float64) float64	
}

type Add struct{}
func (Add) Execute(a, b float64) float64 {return a + b}

type Subtract struct{}
func (Subtract) Execute(a, b float64) float64 {return a - b}

// Calculator управляет доступными операциями через мапу
type Calculator struct {
	operations map[string]Operation
}

// Register позволяет добавлять новые операции БЕЗ изменения кода калькулятора (SOLID)
func (c Calculator) Register(name string, op Operation) {
	c.operations[name] = op
}

// NewCalculator инициализирует калькулятор базовыми операциями
func NewCalculator() *Calculator {
	c := &Calculator{
		operations: make(map[string]Operation),
	}
	c.Register("+", Add{})
	c.Register("-", Subtract{})
	return c
}

// Evaluate вычисляет выражение в обратной польской нотации (например: "3 4 + 2 -")
func (c *Calculator) Evaluate(expr string) (float64, error) {
	tokens := strings.Fields(expr)
	stack := containers.NewStack[float64]()

	for _, token := range tokens {
		// Если токен — это зарегистрированная операция
		if op, exists := c.operations[token]; exists {
			b, okB := stack.Pop()
			a, okA := stack.Pop()
			if !okA || !okB {
				return 0, errors.New("невалидное выражение: не хватает операндов")
			}

			// Выполняем интерфейсный метод
			result := op.Execute(a, b)
			stack.Push(result)
			continue
		}

		// Если это число — парсим и кладем в стек
		val, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return 0, errors.New("неизвестный токен: " + token)
		}
		stack.Push(val)
	}

	res, ok := stack.Pop()
	if !ok {
		return 0, errors.New("пустое выражение")
	}

	if !stack.IsEmpty() {
		return 0, errors.New("невалидное выражение: лишние числа")
	}

	return res, nil
}
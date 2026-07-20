// math_test.go
package math

import (
        "testing"
        "funcs/utils"
        )

func TestAdd(t *testing.T) {
    cases := []struct {
        a, b, want int
    }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
        {100, -50, 50},
    }
    for _, tc := range cases {
        got := Add(tc.a, tc.b)
        if got != tc.want {
            t.Errorf("Add(%d, %d) = %d, хотели %d", tc.a, tc.b, got, tc.want)
        }
    }
}

func TestMinMaxExternal(t *testing.T) {
    //структура одного тест-кейса
    type testCase struct {
        name            string //имя сценария
        input           []int
        expectedMin     int 
        expectedMax     int
    }

    //Массив (таблица) со всеми сценариями (включая edge cases)
    tests := []testCase{
        {
            name:        "Обычный массив",
            input:       []int{5, 4, 9, 7, 2},
            expectedMin: 2,
            expectedMax: 9,
        },
        {
            name:        "Пустой срез (Edge Case)",
            input:       []int{},
            expectedMin: 0,
            expectedMax: 0,
        },
        {
            name:        "Один элемент (Edge Case)",
            input:       []int{42},
            expectedMin: 42,
            expectedMax: 42,
        },
        {
            name:        "Все элементы одинаковые (Edge Case)",
            input:       []int{7, 7, 7, 7},
            expectedMin: 7,
            expectedMax: 7,
        },
        {
            name:        "Отрицательные числа",
            input:       []int{-5, -1, -10, 0, -3},
            expectedMin: -10,
            expectedMax: 0,
        },
    }

    // Перебираем все тест-кейсы в цикле
    for _, tc := range tests {
        // t.Run позволяет запускать каждый сценарий как отдельный подтест
        t.Run(tc.name, func(t *testing.T) {
            gotMin, gotMax := utils.MinMax(tc.input)

            if gotMin != tc.expectedMin || gotMax != tc.expectedMax {
                t.Errorf("Для %s с входными данными %v:\n"+
                    "ожидали: min=%d, max=%d\n"+
                    "получили: min=%d, max=%d",
                    tc.name, tc.input, tc.expectedMin, tc.expectedMax, gotMin, gotMax)
            }            
        })

    }
}

func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}
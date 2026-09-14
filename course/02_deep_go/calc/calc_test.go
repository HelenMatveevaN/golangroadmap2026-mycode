package calc

import "testing"

/*
go test -bench=. -benchmem

go test -v -cover
go test -gcflags="-m -l"

*/

// 1. Создаем НОВУЮ операцию прямо в тестовом файле
type Multiply struct{}

func (Multiply) Execute(a, b float64) float64 {
	return a * b
}

func TestCalculator_Evaluate(t *testing.T) {
	c := NewCalculator()

	c.Register("*", Multiply{})

	tests := []struct {
		name		string
		expr		string
		want 		float64	
		wantErr 	bool
	}{
		{
			name:    "Базовое сложение и вычитание",
			expr:    "3 4 + 2 -", // (3 + 4) - 2 = 5
			want:    5,
			wantErr: false,
		},
		{
			name:    "Проверка новой операции умножения",
			expr:    "3 4 *", // 3 * 4 = 12
			want:    12,
			wantErr: false,
		},
		{
			name:    "Сложное выражение",
			expr:    "5 1 2 + 4 * + 3 -", // 5 + ((1 + 2) * 4) - 3 = 14
			want:    14,
			wantErr: false,
		},
		{
			name:    "Ошибка: не хватает операндов",
			expr:    "3 +",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Ошибка: неизвестный токен",
			expr:    "3 4 @",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Ошибка: лишние числа",
			expr:    "3 4 5 +",
			want:    0,
			wantErr: true,
		},		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Evaluate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
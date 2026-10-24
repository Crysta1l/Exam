package calculation

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"math"
	"testing"
)

func TestCalculateByFactors(t *testing.T) {
	tests := []struct {
		name                      string
		raw                       int
		p1, p2, coefficient, loss float64
		want                      int
	}{
		{"без потерь", 100, 2, 5, 1, 0, 10},
		{"с потерями", 100, 2, 5, 1, 10, 9},
		{"нецелый результат", 10, 3, 1, 1, 0, 3},
		{"десятичная граница", 1, 0.1, 1, 1, 0, 10},
		{"потери в процентах", 1000, 1, 1, 1, 0.12, 998},
		{"коэффициент типа", 100, 2, 5, 1.25, 0, 8},
		{"нулевое сырье", 0, 1, 1, 1, 0, 0},
		{"мало сырья", 1, 2, 1, 1, 0, 0},
		{"полные потери", 100, 1, 1, 1, 100, 0},
		{"отрицательное сырье", -1, 1, 1, 1, 0, -1},
		{"нулевой параметр", 10, 0, 1, 1, 0, -1},
		{"отрицательный параметр", 10, 1, -1, 1, 0, -1},
		{"нулевой коэффициент", 10, 1, 1, 0, 0, -1},
		{"неверные потери", 10, 1, 1, 1, 101, -1},
		{"отрицательные потери", 10, 1, 1, 1, -1, -1},
		{"NaN", 10, math.NaN(), 1, 1, 0, -1},
		{"Inf", 10, 1, math.Inf(1), 1, 0, -1},
		{"переполнение результата", 100, 1e-300, 1e-300, 1, 0, -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CalculateByFactors(test.raw, test.p1, test.p2, test.coefficient, test.loss); got != test.want {
				t.Fatalf("получено %d, ожидалось %d", got, test.want)
			}
		})
	}
}

type fakeDatabase struct{ fail bool }
type fakeRow struct{ fail bool }

func (db fakeDatabase) QueryRow(_ context.Context, _ string, args ...any) pgx.Row {
	return fakeRow{db.fail || args[0].(int) != 1 || args[1].(int) != 1}
}
func (row fakeRow) Scan(dest ...any) error {
	if row.fail {
		return errors.New("нет такого типа или соединения")
	}
	*dest[0].(*float64) = 1
	*dest[1].(*float64) = 10
	return nil
}

func TestCalculateProductionTypes(t *testing.T) {
	for _, ids := range [][2]int{{0, 1}, {1, 0}, {-1, 1}, {999, 1}, {1, 999}} {
		if got := CalculateProduction(context.Background(), fakeDatabase{}, ids[0], ids[1], 100, 2, 5); got != -1 {
			t.Fatalf("неверные типы %v: получено %d", ids, got)
		}
	}
	if got := CalculateProduction(context.Background(), fakeDatabase{}, 1, 1, 100, 2, 5); got != 9 {
		t.Fatalf("ожидалось 9, получено %d", got)
	}
	if got := CalculateProduction(context.Background(), fakeDatabase{true}, 1, 1, 100, 2, 5); got != -1 {
		t.Fatal("ошибка БД должна давать -1")
	}
}

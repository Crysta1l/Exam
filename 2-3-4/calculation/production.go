package calculation

import (
	"context"
	"github.com/jackc/pgx/v5"
	"math"
)

type Database interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func CalculateProduction(ctx context.Context, db Database, productTypeID, materialTypeID, rawAmount int, parameter1, parameter2 float64) int {
	if productTypeID <= 0 || materialTypeID <= 0 || rawAmount < 0 || !positive(parameter1) || !positive(parameter2) {
		return -1
	}
	var coefficient, lossPercent float64
	err := db.QueryRow(ctx, `SELECT pt.coefficient, mt.loss_percent
		FROM product_types pt CROSS JOIN material_types mt
		WHERE pt.id = $1 AND mt.id = $2`, productTypeID, materialTypeID).Scan(&coefficient, &lossPercent)
	if err != nil {
		return -1
	}
	return CalculateByFactors(rawAmount, parameter1, parameter2, coefficient, lossPercent)
}

func positive(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}

func CalculateByFactors(rawAmount int, parameter1, parameter2, coefficient, lossPercent float64) int {
	if rawAmount < 0 || !positive(parameter1) || !positive(parameter2) || !positive(coefficient) || math.IsNaN(lossPercent) || math.IsInf(lossPercent, 0) || lossPercent < 0 || lossPercent > 100 {
		return -1
	}
	available := float64(rawAmount) * (1 - lossPercent/100)
	perItem := parameter1 * parameter2 * coefficient
	if !positive(perItem) {
		return -1
	}
	result := math.Floor(available / perItem)
	if math.IsNaN(result) || math.IsInf(result, 0) || result >= float64(^uint(0)>>1) {
		return -1
	}
	return int(result)
}

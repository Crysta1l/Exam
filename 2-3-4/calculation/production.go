package calculation

import (
	"context"
	"github.com/jackc/pgx/v5"
	"math"
	"math/big"
	"strconv"
)

type Database interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

// The database stores percentages: 0.12 means 0.12%, not a fraction of 0.12.
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
	// Exact decimal arithmetic prevents losing a whole item at a floating-point boundary.
	decimal := func(value float64) *big.Rat {
		result, _ := new(big.Rat).SetString(strconv.FormatFloat(value, 'f', -1, 64))
		return result
	}
	usable := new(big.Rat).Sub(big.NewRat(1, 1), new(big.Rat).Quo(decimal(lossPercent), big.NewRat(100, 1)))
	usable.Mul(usable, new(big.Rat).SetInt64(int64(rawAmount)))
	perItem := new(big.Rat).Mul(decimal(parameter1), decimal(parameter2))
	perItem.Mul(perItem, decimal(coefficient))
	count := new(big.Rat).Quo(usable, perItem)
	whole := new(big.Int).Quo(count.Num(), count.Denom())
	maxInt := new(big.Int).SetUint64(uint64(^uint(0) >> 1))
	if whole.Cmp(maxInt) > 0 {
		return -1
	}
	return int(whole.Int64())
}

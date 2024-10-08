package bricks

import "github.com/shopspring/decimal"

type MoneyAmount struct {
	Value    decimal.Decimal
	Currency string
}

package bricks

import "github.com/shopspring/decimal"

// MoneyAmount is an amount of money in one currency.
// Value is the decimal amount. Currency is the currency code chosen by the caller.
type MoneyAmount struct {
	Value    decimal.Decimal
	Currency string
}

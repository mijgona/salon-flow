package model

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"

	"github.com/shopspring/decimal"
)

const defaultCurrency = "RUB"

// Money represents a monetary value with currency.
type Money struct {
	amount   decimal.Decimal
	currency string
}

// NewMoney creates a new Money value object with validation.
func NewMoney(amount decimal.Decimal, currency string) (Money, error) {
	if amount.IsNegative() {
		return Money{}, errs.NewErrValueMustBe("amount", "non-negative")
	}
	if currency == "" {
		currency = defaultCurrency
	}
	return Money{amount: amount, currency: currency}, nil
}

// MustNewMoney creates a new Money or panics.
func MustNewMoney(amount decimal.Decimal, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// NewMoneyRUB creates a Money value in RUB.
func NewMoneyRUB(amount decimal.Decimal) (Money, error) {
	return NewMoney(amount, defaultCurrency)
}

// Zero returns a zero Money value in the default currency.
func ZeroMoney() Money {
	return Money{amount: decimal.Zero, currency: defaultCurrency}
}

// Amount returns the monetary amount.
func (m Money) Amount() decimal.Decimal { return m.amount }

// Currency returns the currency code.
func (m Money) Currency() string { return m.currency }

// Add adds two Money values. They must share the same currency.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errs.NewErrValueMustBe("currency", "the same for addition")
	}
	return Money{amount: m.amount.Add(other.amount), currency: m.currency}, nil
}

// Subtract subtracts another Money value. Result must not be negative.
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errs.NewErrValueMustBe("currency", "the same for subtraction")
	}
	result := m.amount.Sub(other.amount)
	if result.IsNegative() {
		return Money{}, errs.NewErrValueMustBe("result", "non-negative")
	}
	return Money{amount: result, currency: m.currency}, nil
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool { return m.amount.IsZero() }

// IsPositive returns true if the amount is positive.
func (m Money) IsPositive() bool { return m.amount.IsPositive() }

// Equal checks value equality.
func (m Money) Equal(other Money) bool {
	return m.amount.Equal(other.amount) && m.currency == other.currency
}

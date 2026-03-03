package model

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"

	"github.com/shopspring/decimal"
)

// Discount represents a percent-based discount (0–100).
type Discount struct {
	percent int
}

// NewDiscount creates a Discount value object.
func NewDiscount(percent int) (Discount, error) {
	if percent < 0 || percent > 100 {
		return Discount{}, errs.NewErrValueMustBe("discount", "between 0 and 100")
	}
	return Discount{percent: percent}, nil
}

// MustNewDiscount creates a Discount or panics.
func MustNewDiscount(percent int) Discount {
	d, err := NewDiscount(percent)
	if err != nil {
		panic(err)
	}
	return d
}

// ZeroDiscount returns a 0% discount.
func ZeroDiscount() Discount {
	return Discount{percent: 0}
}

// Percent returns the discount percentage.
func (d Discount) Percent() int { return d.percent }

// Apply applies the discount to a Money value.
func (d Discount) Apply(m Money) Money {
	if d.percent == 0 {
		return m
	}
	factor := decimal.NewFromInt(int64(100 - d.percent)).Div(decimal.NewFromInt(100))
	return Money{amount: m.Amount().Mul(factor), currency: m.Currency()}
}

// IsZero returns true if the discount is 0%.
func (d Discount) IsZero() bool { return d.percent == 0 }

// Equal checks value equality.
func (d Discount) Equal(other Discount) bool { return d.percent == other.percent }

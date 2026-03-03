package loyalty

import "github.com/mijgona/salon-crm/internal/pkg/errs"

// Points represents a loyalty points value.
type Points struct {
	value int
}

// NewPoints creates a Points value object with validation.
func NewPoints(value int) (Points, error) {
	if value < 0 {
		return Points{}, errs.NewErrValueMustBe("points", "non-negative")
	}
	return Points{value: value}, nil
}

// MustNewPoints creates a Points or panics.
func MustNewPoints(value int) Points {
	p, err := NewPoints(value)
	if err != nil {
		panic(err)
	}
	return p
}

// ZeroPoints returns a zero Points value.
func ZeroPoints() Points { return Points{value: 0} }

// Value returns the int value.
func (p Points) Value() int { return p.value }

// Add adds two Points values.
func (p Points) Add(other Points) Points {
	return Points{value: p.value + other.value}
}

// Subtract subtracts Points. Returns error if result would be negative.
func (p Points) Subtract(other Points) (Points, error) {
	result := p.value - other.value
	if result < 0 {
		return Points{}, errs.NewErrValueMustBe("points balance", "sufficient for redemption")
	}
	return Points{value: result}, nil
}

// IsZero returns true if points are zero.
func (p Points) IsZero() bool { return p.value == 0 }

// GreaterThanOrEqual checks if p >= other.
func (p Points) GreaterThanOrEqual(other Points) bool { return p.value >= other.value }

// Equal checks value equality.
func (p Points) Equal(other Points) bool { return p.value == other.value }

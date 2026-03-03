package client

import "github.com/mijgona/salon-crm/internal/pkg/errs"

// Severity represents allergy severity levels.
type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// Allergy represents a client's allergy information.
type Allergy struct {
	substance string
	severity  Severity
}

// NewAllergy creates an Allergy value object with validation.
func NewAllergy(substance string, severity Severity) (Allergy, error) {
	if substance == "" {
		return Allergy{}, errs.NewErrValueRequired("substance")
	}
	if severity == "" {
		severity = SeverityLow
	}
	return Allergy{substance: substance, severity: severity}, nil
}

// MustNewAllergy creates an Allergy or panics.
func MustNewAllergy(substance string, severity Severity) Allergy {
	a, err := NewAllergy(substance, severity)
	if err != nil {
		panic(err)
	}
	return a
}

// Substance returns the allergenic substance.
func (a Allergy) Substance() string { return a.substance }

// Severity returns the severity level.
func (a Allergy) AllergyLevel() Severity { return a.severity }

// Equal checks if two allergies are for the same substance (for deduplication).
func (a Allergy) Equal(other Allergy) bool { return a.substance == other.substance }

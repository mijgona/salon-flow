package loyalty

import (
	"time"

	"github.com/google/uuid"
)

// ReferralStatus represents the status of a referral.
type ReferralStatus string

const (
	ReferralStatusPending   ReferralStatus = "pending"
	ReferralStatusCompleted ReferralStatus = "completed"
)

// Referral is an entity representing a client referral.
type Referral struct {
	id               uuid.UUID
	referredClientID uuid.UUID
	status           ReferralStatus
	bonusEarned      int
	createdAt        time.Time
}

// NewReferral creates a new Referral entity.
func NewReferral(referredClientID uuid.UUID) Referral {
	return Referral{
		id:               uuid.New(),
		referredClientID: referredClientID,
		status:           ReferralStatusPending,
		bonusEarned:      0,
		createdAt:        time.Now(),
	}
}

// RestoreReferral creates a Referral from persisted data.
func RestoreReferral(id, referredClientID uuid.UUID, status ReferralStatus, bonusEarned int, createdAt time.Time) Referral {
	return Referral{
		id:               id,
		referredClientID: referredClientID,
		status:           status,
		bonusEarned:      bonusEarned,
		createdAt:        createdAt,
	}
}

// Complete marks the referral as completed with the earned bonus.
func (r *Referral) Complete(bonus int) {
	r.status = ReferralStatusCompleted
	r.bonusEarned = bonus
}

// Getters
func (r Referral) ID() uuid.UUID               { return r.id }
func (r Referral) ReferredClientID() uuid.UUID { return r.referredClientID }
func (r Referral) Status() ReferralStatus      { return r.status }
func (r Referral) BonusEarned() int            { return r.bonusEarned }
func (r Referral) CreatedAt() time.Time        { return r.createdAt }

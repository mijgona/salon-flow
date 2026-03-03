package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"time"

	"github.com/google/uuid"
)

// MasterScheduleRepository defines operations for persisting MasterSchedule aggregates.
type MasterScheduleRepository interface {
	Add(ctx context.Context, tx interface{}, ms *scheduling.MasterSchedule) error
	Update(ctx context.Context, tx interface{}, ms *scheduling.MasterSchedule) error
	GetByMasterAndDate(ctx context.Context, tx interface{}, masterID uuid.UUID, date time.Time) (*scheduling.MasterSchedule, error)
}

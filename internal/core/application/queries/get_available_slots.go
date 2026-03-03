package queries

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

// GetAvailableSlotsQuery holds data for finding available time slots.
type GetAvailableSlotsQuery struct {
	MasterID        uuid.UUID
	Date            time.Time
	ServiceDuration time.Duration
}

// GetAvailableSlotsResult is the result of the GetAvailableSlots query.
type GetAvailableSlotsResult struct {
	Slots []scheduling.TimeSlot
}

// GetAvailableSlotsHandler handles the get available slots query.
type GetAvailableSlotsHandler struct {
	scheduleRepo ports.MasterScheduleRepository
}

// NewGetAvailableSlotsHandler creates a new handler.
func NewGetAvailableSlotsHandler(scheduleRepo ports.MasterScheduleRepository) *GetAvailableSlotsHandler {
	return &GetAvailableSlotsHandler{scheduleRepo: scheduleRepo}
}

// Handle executes the get available slots query.
func (h *GetAvailableSlotsHandler) Handle(ctx context.Context, q GetAvailableSlotsQuery) (*GetAvailableSlotsResult, error) {
	schedule, err := h.scheduleRepo.GetByMasterAndDate(ctx, nil, q.MasterID, q.Date)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return &GetAvailableSlotsResult{Slots: []scheduling.TimeSlot{}}, nil
	}

	slots := schedule.GetAvailableSlots(q.ServiceDuration)
	return &GetAvailableSlotsResult{Slots: slots}, nil
}

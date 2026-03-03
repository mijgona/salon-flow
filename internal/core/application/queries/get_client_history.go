package queries

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// GetClientHistoryQuery holds data for retrieving a client's appointment history.
type GetClientHistoryQuery struct {
	ClientID uuid.UUID
}

// GetClientHistoryResult is the result of the GetClientHistory query.
type GetClientHistoryResult struct {
	Appointments []*scheduling.Appointment
}

// GetClientHistoryHandler handles the get client history query.
type GetClientHistoryHandler struct {
	appointmentRepo ports.AppointmentRepository
}

// NewGetClientHistoryHandler creates a new handler.
func NewGetClientHistoryHandler(appointmentRepo ports.AppointmentRepository) *GetClientHistoryHandler {
	return &GetClientHistoryHandler{appointmentRepo: appointmentRepo}
}

// Handle executes the get client history query.
func (h *GetClientHistoryHandler) Handle(ctx context.Context, q GetClientHistoryQuery) (*GetClientHistoryResult, error) {
	appointments, err := h.appointmentRepo.FindByClientID(ctx, nil, q.ClientID)
	if err != nil {
		return nil, err
	}
	return &GetClientHistoryResult{Appointments: appointments}, nil
}

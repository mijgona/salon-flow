package http

import (
	"net/http"
	"time"

	"github.com/mijgona/salon-crm/internal/core/application/commands"
	"github.com/mijgona/salon-crm/internal/core/application/queries"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AppointmentHandler handles HTTP requests for appointment management.
type AppointmentHandler struct {
	bookAppointment     *commands.BookAppointmentHandler
	cancelAppointment   *commands.CancelAppointmentHandler
	completeAppointment *commands.CompleteAppointmentHandler
	getAvailableSlots   *queries.GetAvailableSlotsHandler
}

// NewAppointmentHandler creates a new AppointmentHandler.
func NewAppointmentHandler(
	bookAppointment *commands.BookAppointmentHandler,
	cancelAppointment *commands.CancelAppointmentHandler,
	completeAppointment *commands.CompleteAppointmentHandler,
	getAvailableSlots *queries.GetAvailableSlotsHandler,
) *AppointmentHandler {
	return &AppointmentHandler{
		bookAppointment:     bookAppointment,
		cancelAppointment:   cancelAppointment,
		completeAppointment: completeAppointment,
		getAvailableSlots:   getAvailableSlots,
	}
}

// Register registers routes on the Echo instance.
func (h *AppointmentHandler) Register(e *echo.Echo) {
	g := e.Group("/api/v1/appointments")
	g.POST("", h.BookAppointment)
	g.POST("/:id/cancel", h.CancelAppointment)
	g.POST("/:id/complete", h.CompleteAppointment)
	g.GET("/available-slots", h.GetAvailableSlots)
}

// BookAppointmentRequest is the HTTP request body for booking an appointment.
type BookAppointmentRequest struct {
	TenantID  string `json:"tenant_id"`
	ClientID  string `json:"client_id"`
	MasterID  string `json:"master_id"`
	SalonID   string `json:"salon_id"`
	ServiceID string `json:"service_id"`
	StartTime string `json:"start_time"`
	Comment   string `json:"comment,omitempty"`
	Source    string `json:"source"`
}

// BookAppointment handles POST /api/v1/appointments
func (h *AppointmentHandler) BookAppointment(c echo.Context) error {
	var req BookAppointmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, _ := uuid.Parse(req.TenantID)
	clientID, _ := uuid.Parse(req.ClientID)
	masterID, _ := uuid.Parse(req.MasterID)
	salonID, _ := uuid.Parse(req.SalonID)
	serviceID, _ := uuid.Parse(req.ServiceID)
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid start_time format"})
	}

	cmd := commands.BookAppointmentCommand{
		TenantID:  tenantID,
		ClientID:  clientID,
		MasterID:  masterID,
		SalonID:   salonID,
		ServiceID: serviceID,
		StartTime: startTime,
		Comment:   req.Comment,
		Source:    req.Source,
	}

	appointmentID, err := h.bookAppointment.Handle(c.Request().Context(), cmd)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": appointmentID.String()})
}

// CancelAppointment handles POST /api/v1/appointments/:id/cancel
func (h *AppointmentHandler) CancelAppointment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.Bind(&req)

	cmd := commands.CancelAppointmentCommand{
		AppointmentID: id,
		Reason:        req.Reason,
	}

	if err := h.cancelAppointment.Handle(c.Request().Context(), cmd); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// CompleteAppointment handles POST /api/v1/appointments/:id/complete
func (h *AppointmentHandler) CompleteAppointment(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	cmd := commands.CompleteAppointmentCommand{AppointmentID: id}

	if err := h.completeAppointment.Handle(c.Request().Context(), cmd); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAvailableSlots handles GET /api/v1/appointments/available-slots
func (h *AppointmentHandler) GetAvailableSlots(c echo.Context) error {
	masterID, err := uuid.Parse(c.QueryParam("master_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid master_id"})
	}

	date, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format (use YYYY-MM-DD)"})
	}

	durationMin := 60 // default 60 minutes
	duration := time.Duration(durationMin) * time.Minute

	result, err := h.getAvailableSlots.Handle(c.Request().Context(), queries.GetAvailableSlotsQuery{
		MasterID:        masterID,
		Date:            date,
		ServiceDuration: duration,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	slots := make([]map[string]string, 0, len(result.Slots))
	for _, s := range result.Slots {
		slots = append(slots, map[string]string{
			"start_time": s.StartTime().Format(time.RFC3339),
			"end_time":   s.EndTime().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, slots)
}

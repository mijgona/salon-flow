package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mijgona/salon-crm/internal/core/application/queries"
)

// CalendarHandler handles HTTP requests for the appointment calendar.
type CalendarHandler struct {
	getCalendar *queries.GetCalendarHandler
}

// NewCalendarHandler creates a new CalendarHandler.
func NewCalendarHandler(getCalendar *queries.GetCalendarHandler) *CalendarHandler {
	return &CalendarHandler{getCalendar: getCalendar}
}

// Register registers calendar routes on the Echo instance.
func (h *CalendarHandler) Register(e *echo.Echo) {
	g := e.Group("/api/v1/calendar")
	g.GET("", h.GetCalendar)
	g.GET("/master/:master_id", h.GetCalendarByMaster)
	g.GET("/salon/:salon_id", h.GetCalendarBySalon)
}

// GetCalendar handles GET /api/v1/calendar?tenant_id=&from=&to=
func (h *CalendarHandler) GetCalendar(c echo.Context) error {
	tenantID, err := uuid.Parse(c.QueryParam("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or missing tenant_id"})
	}

	from, to, err := parseDateRange(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.getCalendar.Handle(c.Request().Context(), queries.GetCalendarQuery{
		TenantID: tenantID,
		From:     from,
		To:       to,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetCalendarByMaster handles GET /api/v1/calendar/master/:master_id?from=&to=
func (h *CalendarHandler) GetCalendarByMaster(c echo.Context) error {
	masterID, err := uuid.Parse(c.Param("master_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid master_id"})
	}

	from, to, err := parseDateRange(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.getCalendar.Handle(c.Request().Context(), queries.GetCalendarQuery{
		MasterID: masterID,
		From:     from,
		To:       to,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetCalendarBySalon handles GET /api/v1/calendar/salon/:salon_id?from=&to=
func (h *CalendarHandler) GetCalendarBySalon(c echo.Context) error {
	salonID, err := uuid.Parse(c.Param("salon_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid salon_id"})
	}

	from, to, err := parseDateRange(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.getCalendar.Handle(c.Request().Context(), queries.GetCalendarQuery{
		SalonID: salonID,
		From:    from,
		To:      to,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// parseDateRange extracts from/to dates from query parameters.
// Defaults: from = today, to = from + 7 days.
func parseDateRange(c echo.Context) (time.Time, time.Time, error) {
	var from, to time.Time
	var err error

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	if fromStr == "" {
		from = time.Now().Truncate(24 * time.Hour)
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'from' date (use YYYY-MM-DD)")
		}
	}

	if toStr == "" {
		to = from.AddDate(0, 0, 7)
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'to' date (use YYYY-MM-DD)")
		}
	}

	if !to.After(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("'to' must be after 'from'")
	}

	return from, to, nil
}

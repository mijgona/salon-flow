package http

import (
	"net/http"

	"github.com/mijgona/salon-crm/internal/core/application/commands"
	"github.com/mijgona/salon-crm/internal/core/application/queries"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ClientHandler handles HTTP requests for client management.
type ClientHandler struct {
	registerClient *commands.RegisterClientHandler
	updateProfile  *commands.UpdateClientProfileHandler
	getClient      *queries.GetClientHandler
	getHistory     *queries.GetClientHistoryHandler
}

// NewClientHandler creates a new ClientHandler.
func NewClientHandler(
	registerClient *commands.RegisterClientHandler,
	updateProfile *commands.UpdateClientProfileHandler,
	getClient *queries.GetClientHandler,
	getHistory *queries.GetClientHistoryHandler,
) *ClientHandler {
	return &ClientHandler{
		registerClient: registerClient,
		updateProfile:  updateProfile,
		getClient:      getClient,
		getHistory:     getHistory,
	}
}

// Register registers routes on the Echo instance.
func (h *ClientHandler) Register(e *echo.Echo) {
	g := e.Group("/api/v1/clients")
	g.POST("", h.RegisterClient)
	g.GET("/:id", h.GetClient)
	g.PUT("/:id", h.UpdateProfile)
	g.GET("/:id/history", h.GetClientHistory)
}

// RegisterClientRequest is the HTTP request body for registering a client.
type RegisterClientRequest struct {
	TenantID           string `json:"tenant_id"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Source             string `json:"source"`
	ReferredByClientID string `json:"referred_by_client_id,omitempty"`
}

// RegisterClient handles POST /api/v1/clients
func (h *ClientHandler) RegisterClient(c echo.Context) error {
	var req RegisterClientRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
	}

	var referredBy uuid.UUID
	if req.ReferredByClientID != "" {
		referredBy, _ = uuid.Parse(req.ReferredByClientID)
	}

	cmd := commands.RegisterClientCommand{
		TenantID:           tenantID,
		Phone:              req.Phone,
		Email:              req.Email,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Source:             req.Source,
		ReferredByClientID: referredBy,
	}

	clientID, err := h.registerClient.Handle(c.Request().Context(), cmd)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": clientID.String()})
}

// GetClient handles GET /api/v1/clients/:id
func (h *ClientHandler) GetClient(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	result, err := h.getClient.Handle(c.Request().Context(), queries.GetClientQuery{ClientID: id})
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         result.Client.ID().String(),
		"first_name": result.Client.ContactInfo().FirstName(),
		"last_name":  result.Client.ContactInfo().LastName(),
		"phone":      result.Client.ContactInfo().Phone().String(),
		"email":      result.Client.ContactInfo().Email(),
		"source":     result.Client.Source().String(),
		"visits":     result.Client.TotalVisits(),
	})
}

// UpdateProfile handles PUT /api/v1/clients/:id
func (h *ClientHandler) UpdateProfile(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req struct {
		Phone             string   `json:"phone"`
		Email             string   `json:"email"`
		FirstName         string   `json:"first_name"`
		LastName          string   `json:"last_name"`
		Birthday          string   `json:"birthday,omitempty"`
		PreferredMasterID string   `json:"preferred_master_id,omitempty"`
		FavoriteServices  []string `json:"favorite_services,omitempty"`
		Channel           string   `json:"channel,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var preferredMaster uuid.UUID
	if req.PreferredMasterID != "" {
		preferredMaster, _ = uuid.Parse(req.PreferredMasterID)
	}

	var favoriteServices []uuid.UUID
	for _, s := range req.FavoriteServices {
		if uid, err := uuid.Parse(s); err == nil {
			favoriteServices = append(favoriteServices, uid)
		}
	}

	cmd := commands.UpdateClientProfileCommand{
		ClientID:          id,
		Phone:             req.Phone,
		Email:             req.Email,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Birthday:          req.Birthday,
		PreferredMasterID: preferredMaster,
		FavoriteServices:  favoriteServices,
		Channel:           req.Channel,
	}

	if err := h.updateProfile.Handle(c.Request().Context(), cmd); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetClientHistory handles GET /api/v1/clients/:id/history
func (h *ClientHandler) GetClientHistory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	result, err := h.getHistory.Handle(c.Request().Context(), queries.GetClientHistoryQuery{ClientID: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	items := make([]map[string]interface{}, 0, len(result.Appointments))
	for _, a := range result.Appointments {
		items = append(items, map[string]interface{}{
			"id":           a.ID().String(),
			"service_name": a.ServiceInfo().Name(),
			"status":       a.Status().String(),
			"start_time":   a.TimeSlot().StartTime().Format("2006-01-02T15:04:05Z07:00"),
			"price":        a.Price().Amount().String(),
		})
	}

	return c.JSON(http.StatusOK, items)
}

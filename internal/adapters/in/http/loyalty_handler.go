package http

import (
	"github.com/mijgona/salon-crm/internal/core/application/commands"
	"github.com/mijgona/salon-crm/internal/core/application/queries"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// LoyaltyHandler handles HTTP requests for loyalty management.
type LoyaltyHandler struct {
	earnPoints        *commands.EarnPointsHandler
	getLoyaltyAccount *queries.GetLoyaltyAccountHandler
}

// NewLoyaltyHandler creates a new LoyaltyHandler.
func NewLoyaltyHandler(
	earnPoints *commands.EarnPointsHandler,
	getLoyaltyAccount *queries.GetLoyaltyAccountHandler,
) *LoyaltyHandler {
	return &LoyaltyHandler{
		earnPoints:        earnPoints,
		getLoyaltyAccount: getLoyaltyAccount,
	}
}

// Register registers routes on the Echo instance.
func (h *LoyaltyHandler) Register(e *echo.Echo) {
	g := e.Group("/api/v1/loyalty")
	g.GET("/:client_id", h.GetLoyaltyAccount)
	g.POST("/:client_id/earn", h.EarnPoints)
}

// GetLoyaltyAccount handles GET /api/v1/loyalty/:client_id
func (h *LoyaltyHandler) GetLoyaltyAccount(c echo.Context) error {
	clientID, err := uuid.Parse(c.Param("client_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client_id"})
	}

	result, err := h.getLoyaltyAccount.Handle(c.Request().Context(), queries.GetLoyaltyAccountQuery{ClientID: clientID})
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":              result.Account.ID().String(),
		"client_id":       result.Account.ClientID().String(),
		"tier":            result.Account.Tier().String(),
		"balance":         result.Account.Balance().Value(),
		"lifetime_points": result.Account.LifetimePoints().Value(),
		"discount":        result.Account.Tier().DiscountPercent(),
	})
}

// EarnPoints handles POST /api/v1/loyalty/:client_id/earn
func (h *LoyaltyHandler) EarnPoints(c echo.Context) error {
	clientID, err := uuid.Parse(c.Param("client_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client_id"})
	}

	var req struct {
		Amount          int    `json:"amount"`
		Reason          string `json:"reason"`
		RelatedEntityID string `json:"related_entity_id,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var relatedEntityID uuid.UUID
	if req.RelatedEntityID != "" {
		relatedEntityID, _ = uuid.Parse(req.RelatedEntityID)
	}

	cmd := commands.EarnPointsCommand{
		ClientID:        clientID,
		Amount:          req.Amount,
		Reason:          req.Reason,
		RelatedEntityID: relatedEntityID,
	}

	if err := h.earnPoints.Handle(c.Request().Context(), cmd); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

package payment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/asino-nelson/safiri-logistics/internal/middleware"
	"github.com/asino-nelson/safiri-logistics/internal/order"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group := router.Group("/payments")
	group.Use(authMiddleware)
	group.POST("/loads/:loadID/mpesa-checkout", h.initiateCheckout)
}

func (h *Handler) initiateCheckout(c *gin.Context) {
	var request CheckoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.InitiateCheckout(
		c.Request.Context(),
		middleware.CurrentUserID(c),
		middleware.CurrentUserRole(c),
		c.Param("loadID"),
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrOnlyLoadOwnerCanPay):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, ErrInvalidQuotedAmount):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, order.ErrLoadNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate payment"})
		}
		return
	}

	c.JSON(http.StatusCreated, payment)
}

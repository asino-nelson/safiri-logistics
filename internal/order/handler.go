package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/asino-nelson/safiri-logistics/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	loads := router.Group("/loads")
	loads.Use(authMiddleware)
	loads.POST("", h.create)
	loads.GET("", h.list)
}

func (h *Handler) create(c *gin.Context) {
	var request CreateLoadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	load, err := h.service.CreateLoad(c.Request.Context(), middleware.CurrentUserID(c), middleware.CurrentUserRole(c), request)
	if err != nil {
		if errors.Is(err, ErrOnlyCustomersCanPostLoads) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create load"})
		return
	}

	c.JSON(http.StatusCreated, load)
}

func (h *Handler) list(c *gin.Context) {
	loads, err := h.service.ListLoads(c.Request.Context(), middleware.CurrentUserID(c), middleware.CurrentUserRole(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list loads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": loads})
}

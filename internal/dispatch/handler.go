package dispatch

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
	group := router.Group("/dispatch")
	group.Use(authMiddleware)
	group.POST("/optimize", h.optimize)
}

func (h *Handler) optimize(c *gin.Context) {
	var request OptimizeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.OptimizeAssignments(c.Request.Context(), middleware.CurrentUserRole(c), request)
	if err != nil {
		if errors.Is(err, ErrOnlyAdminsCanOptimize) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to optimize dispatch"})
		return
	}

	c.JSON(http.StatusOK, result)
}

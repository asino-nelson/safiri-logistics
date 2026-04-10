package driver

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
	drivers := router.Group("/drivers")
	drivers.Use(authMiddleware)
	drivers.POST("/kyc", h.submitKYC)
	drivers.GET("/kyc/me", h.getMyProfile)
	drivers.PATCH("/:userID/kyc/review", h.reviewKYC)
}

func (h *Handler) submitKYC(c *gin.Context) {
	var request SubmitKYCRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.SubmitKYC(c.Request.Context(), middleware.CurrentUserID(c), middleware.CurrentUserRole(c), request)
	if err != nil {
		if errors.Is(err, ErrOnlyDriversCanSubmitKYC) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit kyc"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *Handler) getMyProfile(c *gin.Context) {
	profile, err := h.service.GetProfile(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch kyc profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) reviewKYC(c *gin.Context) {
	var request ReviewKYCRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.ReviewKYC(c.Request.Context(), middleware.CurrentUserRole(c), c.Param("userID"), request)
	if err != nil {
		switch {
		case errors.Is(err, ErrOnlyAdminsCanReviewKYC):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, ErrProfileNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrInvalidKYCStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to review kyc"})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

type Handler struct {
	service *Service
	users   *user.Service
}

func NewHandler(service *Service, users *user.Service) *Handler {
	return &Handler{service: service, users: users}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
	router.GET("/me", authMiddleware, h.me)
}

func (h *Handler) register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(c.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, ErrInvalidRole):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) me(c *gin.Context) {
	account, err := h.users.GetByID(c.Request.Context(), c.GetString("currentUserID"))
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, UserPayload{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
		Role:  account.Role,
	})
}

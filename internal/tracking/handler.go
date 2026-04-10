package tracking

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/asino-nelson/safiri-logistics/internal/middleware"
	"github.com/asino-nelson/safiri-logistics/internal/order"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group := router.Group("/tracking")
	group.Use(authMiddleware)
	group.POST("/loads/:loadID/events", h.createEvent)
	group.GET("/ws", h.websocket)
}

func (h *Handler) createEvent(c *gin.Context) {
	var request CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(
		c.Request.Context(),
		middleware.CurrentUserID(c),
		middleware.CurrentUserRole(c),
		c.Param("loadID"),
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrOnlyAssignedDriverCanTrack):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, order.ErrLoadNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tracking event"})
		}
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *Handler) websocket(c *gin.Context) {
	loadID := c.Query("load_id")
	if loadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "load_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	stream := h.service.Subscribe(loadID)
	defer h.service.Unsubscribe(loadID, stream)

	for {
		select {
		case event, ok := <-stream:
			if !ok {
				return
			}

			if err := conn.WriteJSON(event); err != nil {
				return
			}
		case <-time.After(30 * time.Second):
			if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second)); err != nil {
				return
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}

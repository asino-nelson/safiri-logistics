package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/asino-nelson/safiri-logistics/internal/auth"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tokens := auth.NewJWTManager("secret", time.Hour)
	token, err := tokens.Generate("user-1", "customer")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	router := gin.New()
	router.GET("/protected", Authenticate(tokens), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": CurrentUserID(c),
			"role":    CurrentUserRole(c),
		})
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if body == "" || !containsAll(body, "user-1", "customer") {
		t.Fatalf("unexpected response body: %s", body)
	}
}

func containsAll(value string, expected ...string) bool {
	for _, item := range expected {
		if !strings.Contains(value, item) {
			return false
		}
	}

	return true
}

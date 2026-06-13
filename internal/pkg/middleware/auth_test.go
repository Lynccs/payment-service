package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appjwt "github.com/Lynccs/payment-service/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

func init() {
	gin.SetMode(gin.TestMode)
}

func newAuthRouter() *gin.Engine {
	r := gin.New()
	r.GET("/protected", AuthRequired(testSecret), func(c *gin.Context) {
		id, _ := c.Get(UserIDKey)
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})
	return r
}

func TestAuthRequired_NoCookie(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	newAuthRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid.token.value"})

	newAuthRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	token, err := appjwt.GenerateToken(1, testSecret, -time.Second)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	newAuthRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_WrongSecret(t *testing.T) {
	token, err := appjwt.GenerateToken(1, "other-secret", time.Hour)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	newAuthRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_ValidToken(t *testing.T) {
	token, err := appjwt.GenerateToken(99, testSecret, time.Hour)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	newAuthRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "99")
}

func TestAuthRequired_SetsUserIDInContext(t *testing.T) {
	expectedID := 7
	token, err := appjwt.GenerateToken(expectedID, testSecret, time.Hour)
	require.NoError(t, err)

	var capturedID any
	r := gin.New()
	r.GET("/check", AuthRequired(testSecret), func(c *gin.Context) {
		capturedID, _ = c.Get(UserIDKey)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedID, capturedID)
}

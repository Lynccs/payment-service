package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestGenerateAndParseToken(t *testing.T) {
	userID := 42

	token, err := GenerateToken(userID, testSecret, time.Hour)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedID, err := ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	token, err := GenerateToken(1, testSecret, -time.Second)
	require.NoError(t, err)

	_, err = ParseToken(token, testSecret)
	assert.Error(t, err)
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(1, testSecret, time.Hour)
	require.NoError(t, err)

	_, err = ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_InvalidString(t *testing.T) {
	_, err := ParseToken("not.a.token", testSecret)
	assert.Error(t, err)
}

func TestParseToken_WrongSigningMethod(t *testing.T) {
	// "none" — не є HMAC, сервер має відхилити
	claims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token, err := raw.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = ParseToken(token, testSecret)
	assert.Error(t, err, "токен підписаний не HMAC має відхилятись")
}

func TestGenerateToken_DifferentUserIDs(t *testing.T) {
	t1, _ := GenerateToken(1, testSecret, time.Hour)
	t2, _ := GenerateToken(2, testSecret, time.Hour)

	assert.NotEqual(t, t1, t2)

	id1, err := ParseToken(t1, testSecret)
	require.NoError(t, err)
	id2, err := ParseToken(t2, testSecret)
	require.NoError(t, err)

	assert.Equal(t, 1, id1)
	assert.Equal(t, 2, id2)
}

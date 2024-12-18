package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTAuthToken(t *testing.T) {
	nowTime := time.Now()
	jwtExpiryDuration := time.Minute * 60 * 24 * 365 * 10
	var jwtSecretKey = []byte("capoo_is_cute")

	// encode JWT
	userID := "00000000-0000-0000-0000-000000000001"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "https://our.domain.com",
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(nowTime.Add(jwtExpiryDuration)),
		ID:        uuid.New().String(),
	})
	signedKey, err := token.SignedString(jwtSecretKey)
	assert.NoError(t, err)

	fmt.Printf("user ID: %s\n", userID)
	fmt.Printf("expiry time: %v\n", nowTime.Add(jwtExpiryDuration).Format(time.RFC3339))
	fmt.Printf("signedKey: %s\n", signedKey)

}

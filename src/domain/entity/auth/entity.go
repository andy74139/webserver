package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	CreateToken(ctx context.Context, userID uuid.UUID) (string, error)
	ParseToken(ctx context.Context, jwtToken string) (*jwt.RegisteredClaims, error)
	ParseAndVerifyToken(ctx context.Context, jwtToken string) (*jwt.RegisteredClaims, bool, error)
	RevokeToken(ctx context.Context, jwtID string, jwtExpiryTime time.Time) error
}

type Repository interface {
	SetRevoked(ctx context.Context, jwtID string, expiryDuration time.Duration) error
	IsRevoked(ctx context.Context, jwtID string) (bool, error)
}

var (
	ErrRevokedToken = errors.New("revoked token")
)

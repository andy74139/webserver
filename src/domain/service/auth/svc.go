package auth_svc

import (
	"context"
	"fmt"
	"time"

	"github.com/andy74139/webserver/src/domain/entity/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	jwtIssuer                 = "https://my.domain.com"
	jwtExpiryDuration         = time.Hour * 24 * 7
	jwtSuggestRefreshDuration = time.Hour * 24 * 3
)

var jwtSecretKey = []byte("secret_key")
var now = time.Now

type service struct {
	repo auth.Repository
}

func New(repo auth.Repository) auth.Service {
	return &service{repo: repo}
}

func (s *service) CreateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    jwtIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now()),
		ExpiresAt: jwt.NewNumericDate(now().Add(jwtExpiryDuration)),
		ID:        uuid.New().String(),
	})
	signedKey, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("SingedString error: %w", err)
	}

	return signedKey, nil
}

func (s *service) ParseToken(ctx context.Context, jwtToken string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) { return jwtSecretKey, nil })
	if err != nil {
		return nil, fmt.Errorf("jwt.Parse error: %w", err)
	}
	return claims, nil
}

func (s *service) ParseAndVerifyToken(ctx context.Context, jwtToken string) (*jwt.RegisteredClaims, bool, error) {
	claim, err := s.ParseToken(ctx, jwtToken)
	if err != nil {
		return nil, false, err
	}

	// verify
	issuer, err := claim.GetIssuer()
	if err != nil {
		return nil, false, fmt.Errorf("GetIssuer error: %w", err)
	} else if issuer != jwtIssuer {
		return nil, false, fmt.Errorf("invalid issuer: %s", issuer)
	}

	if _, err := claim.GetSubject(); err != nil {
		return nil, false, fmt.Errorf("GetSubject error: %w", err)
	}
	expiryDur, err := claim.GetExpirationTime()
	if err != nil {
		return nil, false, fmt.Errorf("GetExpirationTime error: %w", err)
	}

	if isRevoked, err := s.repo.IsRevoked(ctx, claim.ID); err != nil {
		return nil, false, fmt.Errorf("repo.IsRevoked error: %w", err)
	} else if isRevoked {
		return nil, false, auth.ErrRevokedToken
	}

	isSuggestRefresh := expiryDur.Sub(now()) <= jwtSuggestRefreshDuration
	return claim, isSuggestRefresh, nil
}

func (s *service) RevokeToken(ctx context.Context, jwtID string, jwtExpiryTime time.Time) error {
	expiryDuration := jwtExpiryTime.Sub(now())
	if err := s.repo.SetRevoked(ctx, jwtID, expiryDuration); err != nil {
		return fmt.Errorf("SetRevoke error: %w", err)
	}
	return nil
}

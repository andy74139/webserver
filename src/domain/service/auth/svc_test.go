package auth_svc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (s *AuthSuite) TestJWT_HappyCase() {
	ctx := context.TODO()
	svc := New(nil)
	nowTime := time.Now()

	// override now() used in app
	now = func() time.Time { return nowTime }

	// encode JWT
	userID := uuid.New()
	token, err := svc.CreateToken(ctx, userID)
	s.Require().NoError(err)
	fmt.Println(token)

	// decode JWT
	cliams, err := svc.ParseToken(ctx, token)
	s.Require().NoError(err)

	issuer, _ := cliams.GetIssuer()
	subject, _ := cliams.GetSubject()
	issuedAt, _ := cliams.GetIssuedAt()
	expiresAt, _ := cliams.GetExpirationTime()
	s.Equal(jwtIssuer, issuer)
	s.Equal(userID.String(), subject)
	s.Equal(nowTime.Unix(), issuedAt.Unix())
	s.Equal(nowTime.Add(jwtExpiryDuration).Unix(), expiresAt.Unix())
}

func (s *AuthSuite) TestJWT_TimeExpired() {
	ctx := context.TODO()
	nowTime := time.Now().Add(-jwtExpiryDuration)
	svc := New(nil)

	// override now() used in app
	now = func() time.Time { return nowTime }

	// encode JWT
	userID := uuid.New()
	token, err := svc.CreateToken(ctx, userID)
	s.Require().NoError(err)
	fmt.Println(token)

	// decode JWT
	nowTime = nowTime.Add(jwtExpiryDuration)
	_, err = svc.ParseToken(ctx, token)
	s.Require().ErrorIs(err, jwt.ErrTokenExpired)
}

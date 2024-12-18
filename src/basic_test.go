package src

// The tests are asserted for some basic functions.
// And gives a clear view for developers not knowing them.

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type BasicSuite struct {
	suite.Suite
}

func TestBasicSuite(t *testing.T) {
	// TODO: let it can be run using `make basic`
	t.Skip() // comment it when you want to check basic functions

	suite.Run(t, new(BasicSuite))
}

func (s *BasicSuite) TestJWT_HappyCase() {
	secret := []byte("capoo_is_cute")

	// encode JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "https://our.domain.com",
		Subject:   "01B6C691-EE6E-48F2-B2DE-AFC1DA180EFE",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	signedKey, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println(signedKey)

	// decode JWT
	token2, err := jwt.Parse(signedKey, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		panic(err)
	}
	issuer, _ := token2.Claims.GetIssuer()
	subject, _ := token2.Claims.GetSubject()
	issuedAt, _ := token2.Claims.GetIssuedAt()
	expiresAt, _ := token2.Claims.GetExpirationTime()
	fmt.Println(issuer)
	fmt.Println(subject)
	fmt.Println(issuedAt)
	fmt.Println(expiresAt)
}

func (s *BasicSuite) TestJWT_TimeExpired() {
	secret := []byte("capoo_is_cute")

	// encode JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "https://our.domain.com",
		Subject:   "01B6C691-EE6E-48F2-B2DE-AFC1DA180EFE",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now()),
	})
	signedKey, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println(signedKey)

	// decode JWT
	if _, err := jwt.Parse(signedKey, func(t *jwt.Token) (interface{}, error) { return secret, nil }); errors.Is(err, jwt.ErrTokenExpired) {
		fmt.Println(err)
	} else {
		panic("should have token expired error")
	}
}

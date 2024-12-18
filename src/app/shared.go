package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/andy74139/webserver/src/infra"
)

// defaultRequestHeaders represents the model for header params
// @HeaderParameters defaultRequestHeaders
type defaultRequestHeaders struct {
	Authorization string `json:"Authorization" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL291ci5kb21haW4uY29tIiwic3ViIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwiZXhwIjoyMDQ3MDE1NTQ2LCJqdGkiOiI5Yzc1MGU5OC1jNDU2LTRhYzAtODY1Yy1kM2Q5ZDY4MDE2YjUifQ.X2nk0rfeEzcgthvZXoP2wi6k63JF4kJVg6I3M-Alals"`
}

// verifyAuth is a shared method for endpoint methods
func (a *app) verifyAuth(ctx *gin.Context) (uuid.UUID, *jwt.RegisteredClaims, bool, bool) {
	logger := infra.GetLogger(ctx)

	jwtString := a.getJWTString(ctx)
	claims, isSuggestRefresh, err := a.authSvc.ParseAndVerifyToken(ctx, jwtString)
	if err != nil {
		logger.Debugw("GetSubject error", "error", fmt.Errorf("ParseAndVerifyToken error: %w", err))
		ctx.AbortWithStatus(http.StatusForbidden)
		return uuid.Nil, claims, false, false
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		logger.Debugw("uuid.Parse error", "error", fmt.Errorf("uuid.Parse error: %w", err))
		ctx.AbortWithStatus(http.StatusForbidden)
		return uuid.Nil, claims, false, false
	}
	return userID, claims, isSuggestRefresh, true
}

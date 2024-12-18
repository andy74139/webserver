package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/andy74139/webserver/src/domain/entity/user"
	"github.com/andy74139/webserver/src/infra"
)

type requestLoginByDevice struct {
	Platform string `json:"platform" example:"android" description:"Platform of the device"`
	DeviceID string `json:"device_id" example:"123456" description:"Device ID"`
	Name     string `json:"name,omitempty" example:"User123456" description:"User name"`
}

type responseAuthToken struct {
	AuthToken string `json:"auth" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL291ci5kb21haW4uY29tIiwic3ViIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwiZXhwIjoyMDQ3MDE1NTQ2LCJqdGkiOiI5Yzc1MGU5OC1jNDU2LTRhYzAtODY1Yy1kM2Q5ZDY4MDE2YjUifQ.X2nk0rfeEzcgthvZXoP2wi6k63JF4kJVg6I3M-Alals" description:"Authorization token"`
}

// @Title Login by device
// @Description Login by device, which creates an authorization token. It registers an account if it doesn't exist.
// @Tags authorization
// @Resource authorization
// @Accept json
// @Produce json
// @Param request body requestLoginByDevice true "Login request"
// @Success 200 object responseAuthToken "OK"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /api/v1/account/auth [post]
func (a *app) loginByDevice(ctx *gin.Context) {
	// NOTE: To be convenient to front-end, the API registers an account if the account doesn't exist

	logger := infra.GetLogger(ctx)

	req := &requestLoginByDevice{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Debugw("unknown request body", "request", ctx.Request)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown request body"})
		return
	}

	// get and create(if need) user
	userID, err := a.userSvc.GetIDByDevice(ctx, req.Platform, req.DeviceID)
	if errors.Is(err, user.ErrNotFound) {
		// register it if user is not found
		logger.Debugw("user not exist, register one", "request", req)
		if err := a.userSvc.Create(ctx, req.Platform, req.DeviceID, req.Name); err != nil {
			logger.Errorw("Create error", "req", req, "error", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		userID, err = a.userSvc.GetIDByDevice(ctx, req.Platform, req.DeviceID)
		if err != nil {
			logger.Errorw("GetIDByDevice error", "error", err, "request", req)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		logger.Errorw("GetIDByDevice error", "error", err, "request", req)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create auth token
	token, err := a.authSvc.CreateToken(ctx, userID)
	if err != nil {
		logger.Errorw("CreateToken error", "error", err, "user_id", userID.ID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO: 201 Created
	ctx.JSON(http.StatusOK, &responseAuthToken{AuthToken: token})
}

func (a *app) loginBySSO(ctx *gin.Context) {
	panic("UNDONE")

	type request struct {
		SSOProvider string `json:"sso_provider"`
		SSOToken    string `json:"sso_token"`
	}
	logger := infra.GetLogger(ctx)

	req := &request{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Debugw("unknown request body", "request", ctx.Request)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown request body"})
		return
	}

	//TODO: get account ID
	ssoAccountID := ""

	userID, err := a.userSvc.GetIDBySSO(ctx, req.SSOProvider, ssoAccountID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			logger.Debugw("user not exist", "sso_provider", req.SSOProvider, "sso_account_id", ssoAccountID)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
			return
		}
		logger.Errorw("GetIDBySSO error", "error", err, "request", req)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	token, err := a.authSvc.CreateToken(ctx, userID)
	if err != nil {
		logger.Errorw("CreateToken error", "error", err, "user_id", userID.ID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &responseAuthToken{AuthToken: token})
}

// @Title Refresh authorization token
// @Description Revokes the authorization token, and create a new one.
// @Tags authorization
// @Resource authorization
// @Accept json
// @Produce json
// @Header defaultRequestHeaders
// @Success 200 object responseAuthToken "OK"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /api/v1/account/auth [put]
func (a *app) refreshAuthToken(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)
	userID, claims, _, ok := a.verifyAuth(ctx)
	if !ok {
		return
	}

	token, err := a.authSvc.CreateToken(ctx, userID)
	if err != nil {
		logger.Errorw("CreateToken error", "error", err, "user_id", userID.ID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := a.authSvc.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		logger.Errorw("RevokeToken error", "error", err, "jwt_id", claims.ID)
		// NOTE: still return ok to client
		// TODO: retry revoke token
	}

	ctx.JSON(http.StatusOK, gin.H{"auth": token})
}

// @Title Logout
// @Description Revokes the authorization token.
// @Tags authorization
// @Resource authorization
// @Accept json
// @Produce json
// @Header defaultRequestHeaders
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /api/v1/account/auth [delete]
func (a *app) logout(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)

	_, claims, _, ok := a.verifyAuth(ctx)
	if !ok {
		return
	}

	if err := a.authSvc.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		logger.Errorw("RevokeToken error", "error", err, "claims_id", claims.ID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

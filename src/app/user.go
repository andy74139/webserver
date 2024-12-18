package app

import (
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/andy74139/webserver/src/domain/entity/user"
	"github.com/andy74139/webserver/src/infra"
)

type requestRegister struct {
	Platform string `json:"platform" example:"android" description:"Platform of the device"`
	DeviceID string `json:"device_id" example:"123456" description:"Device ID"`
	Name     string `json:"name,omitempty" example:"User123456" description:"User name"`
}

// @Title Register user
// @Description Register user by device
// @Param  request  body  requestRegister  true  "Request body"
// @Success  201  "Created"
// @Failure  400  "Bad Request"
// @Failure  500  "Internal Server Error"
// @Resource account
// @Route /api/v1/account [post]
func (a *app) registerByDevice(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)

	req := &requestRegister{}
	if err := ctx.BindJSON(req); err != nil {
		// TODO: check can still read body after BindJSON
		if body, err2 := io.ReadAll(ctx.Request.Body); err2 != nil {
			logger.Debugw("failed to read request body", "request", ctx.Request, "BindJSON error", err, "ReadAll error", err2)
		} else {
			logger.Debugw("BindJSON error", "body", body, "error", err)
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown request body"})
		return
	}

	err := a.userSvc.Create(ctx, req.Platform, req.DeviceID, req.Name)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		logger.Errorw("Create error", "req", req, "error", err)
		return
	}

	ctx.Status(http.StatusCreated)
}

type responseGetUserInfo struct {
	ID               string `json:"id" example:"01B6C691-EE6E-48F2-B2DE-AFC1DA180EFE" description:"User ID"`
	Name             string `json:"name" example:"User123456" description:"User name"`
	IsSuggestRefresh bool   `json:"is_suggest_refresh" example:"false" description:"Is authorization token suggest to refresh"`
}

// @Title Get user info
// @Description Get user info
// @Header defaultRequestHeaders
// @Success  200  object  responseGetUserInfo  "OK"
// @Failure  403  "Forbidden"
// @Failure  500  "Internal Server Error"
// @Resource account
// @Route /api/v1/account [get]
func (a *app) getUserInfo(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)
	userID, _, isSuggestRefresh, ok := a.verifyAuth(ctx)
	if !ok {
		return
	}

	user1, err := a.userSvc.Get(ctx, userID)
	if err != nil {
		logger.Errorw("userSvc.Get error", "user_id", userID, "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &responseGetUserInfo{
		ID:               user1.ID.String(),
		Name:             user1.Name,
		IsSuggestRefresh: isSuggestRefresh,
	})
}

func (a *app) updateUserInfo(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)
	userID, _, isSuggestRefresh, ok := a.verifyAuth(ctx)
	if !ok {
		return
	}

	type request struct {
		Name string `json:"name,omitempty"`
	}
	req := &request{}
	if err := ctx.BindJSON(req); err != nil {
		// TODO: check can still read body after BindJSON
		if body, err2 := io.ReadAll(ctx.Request.Body); err2 != nil {
			logger.Debugw("failed to read request body", "request", ctx.Request, "BindJSON error", err, "ReadAll error", err2)
		} else {
			logger.Debugw("BindJSON error", "body", body, "error", err)
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown request body"})
		return
	}

	user1 := &user.User{
		ID:   userID,
		Name: req.Name,
	}
	err := a.userSvc.Update(ctx, user1)
	if err != nil {
		logger.Errorw("userSvc.Get error", "user_id", userID, "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":               user1.Name,
		"is_suggest_refresh": isSuggestRefresh,
	})
}

func (a *app) addSSO(ctx *gin.Context) {
	panic("UNDONE")

	//logger := infra.GetLogger(ctx)
	//userID, _, isSuggestRefresh, ok := a.verifyAuth(ctx)
	//if !ok {
	//	return
	//}
	//
	//type request struct {
	//	Platform string `json:"platform"`
	//	DeviceID string `json:"device_id"`
	//}
	//req := &request{}
	//if err := ctx.BindJSON(req); err != nil {
	//	// TODO: check can still read body after BindJSON
	//	if body, err2 := io.ReadAll(ctx.Request.Body); err2 != nil {
	//		logger.Debugw("failed to read request body", "request", ctx.Request, "BindJSON error", err, "ReadAll error", err2)
	//	} else {
	//		logger.Debugw("BindJSON error", "body", body, "error", err)
	//	}
	//
	//	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown request body"})
	//	return
	//}
	//
	//// TODO: get sso account id
	//ssoAccountID := ""
	//
	//err := a.userSvc.AddSSO(ctx, userID, req.Platform, ssoAccountID)
	//if err != nil {
	//	logger.Errorw("userSvc.AddSSO error", "user_id", userID, "error", err)
	//	ctx.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, gin.H{
	//	"is_suggest_refresh": isSuggestRefresh,
	//})
}

// @Title Delete user
// @Description Remove the account
// @Header defaultRequestHeaders
// @Success  200  "OK"
// @Failure  403  "Forbidden"
// @Failure  500  "Internal Server Error"
// @Resource account
// @Route /api/v1/account [delete]
func (a *app) deleteAccount(ctx *gin.Context) {
	logger := infra.GetLogger(ctx)
	userID, _, _, ok := a.verifyAuth(ctx)
	if !ok {
		return
	}

	if err := a.userSvc.Delete(ctx, userID); err != nil {
		logger.Errorw("userSvc.Delete error", "user_id", userID, "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//TODO: 204 No Content
	ctx.Status(http.StatusOK)
}

func (a *app) getJWTString(ctx *gin.Context) string {
	auth := ctx.Request.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

var runeSet = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func getRandomRequestID() string {
	const length = 8
	if len(runeSet) != 36 {
		panic("ya")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = runeSet[rand.Intn(len(runeSet))]
	}
	return string(result)
}

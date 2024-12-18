package infra

import (
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetGinLogger(loggerName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := GetDefaultLogger().Named(loggerName)
		ctx = ginWithLogger(ctx, logger)
		ctx.Set(loggerKey, logger)

		//ctx.Set("request_id", id)
		//ctx.Writer.Header().Add("request-id", id)

		//logger.With("request_id", id)
		//logger.With("client_ip", ctx.ClientIP())
		//logger.With("http_method", ctx.Request.Method)
		//logger.With("host", ctx.Request.Host)
		//logger.With("url_path", ctx.Request.URL.EscapedPath())

		defer func(t time.Time) {
			//logger.With("latency", strconv.FormatInt(time.Since(t).Milliseconds(), 10))
			//logger.With("status_code", strconv.Itoa(ctx.Writer.Status()))
			logger.Info("request")
		}(time.Now())

		ctx.Next()
	}
}

func ginWithLogger(ctx *gin.Context, logger *zap.SugaredLogger) *gin.Context {
	reqCtx := ctx.Request.Context()
	reqCtx = SetLogger(reqCtx, logger)
	ctx.Request = ctx.Request.WithContext(reqCtx)
	return ctx
}

// PanicCatcher defines a panic catcher handler.
func PanicCatcher(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger := GetLogger(ctx)
			logger.Errorf("panic_catcher: %v\n%s", err, debug.Stack())
		}
	}()
	// Process Request Chain
	ctx.Next()
}

// TODO: prevent long request middeleware

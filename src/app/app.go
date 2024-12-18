package app

// App is an instance of backend server.
// Putting most coherency here, instead of main program.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"

	"github.com/andy74139/webserver/src/config"
	"github.com/andy74139/webserver/src/domain/entity/auth"
	"github.com/andy74139/webserver/src/domain/entity/user"
	"github.com/andy74139/webserver/src/domain/repository/auth"
	"github.com/andy74139/webserver/src/domain/repository/user"
	"github.com/andy74139/webserver/src/domain/service/auth"
	"github.com/andy74139/webserver/src/domain/service/user"
	"github.com/andy74139/webserver/src/infra"
)

// @Version 1.0.0
// @Title WebServer Backend API
// @Description Backend API
// @ContactName Andy
// @ContactURL http://domain.com
// @TermsOfServiceUrl http://someurl.oxox
// @Server http://www.domain.com Server-1
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token

type App interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type app struct {
	server *http.Server

	userSvc user.Service
	authSvc auth.Service
}

func New() App {
	return &app{}
}

func (a *app) Start(ctx context.Context) error {
	// logger
	// TODO: info log file
	writeSyncer, _, err := zap.Open("/app/log/error.log")
	if err != nil {
		panic(fmt.Errorf("zap.Open error: %w", err))
	}
	log := zap.Must(zap.NewDevelopment(zap.ErrorOutput(writeSyncer)))
	logger := log.Sugar()
	logger.Info("Start!")
	defer logger.Sync()

	infra.SetDefaultLogger(logger)
	ctx = infra.SetLogger(ctx, logger)

	a.setDomainServices()

	// HTTP server
	router := a.getRouter(ctx)
	port := config.GetWebServerPort()
	a.server = &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	// Start Running HTTP Server.
	if err := a.server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		return nil
	} else if err != nil {
		return fmt.Errorf("gin ListenAndServe error: %w", err)
	}
	return nil
}

func (a *app) Stop(ctx context.Context) error {
	if a.server == nil {
		return errors.New("app not started")
	}

	return a.server.Shutdown(ctx)
}

func (a *app) getRouter(ctx context.Context) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.ForwardedByClientIP = true
	router.UseRawPath = false
	router.UnescapePathValues = true
	router.Use(infra.PanicCatcher)

	// Account, register account
	accountRouter := router.Group("/api/v1/account")
	accountRouter.POST("/",
		infra.SetGinLogger("account_create"),
		a.registerByDevice,
	)
	accountRouter.GET("/",
		infra.SetGinLogger("account_get_info"),
		a.getUserInfo,
	)
	accountRouter.PUT("/",
		infra.SetGinLogger("account_update_info"),
		a.updateUserInfo,
	)
	accountRouter.PUT("/sso",
		infra.SetGinLogger("account_update_sso"),
		a.addSSO,
	)
	accountRouter.DELETE("/",
		infra.SetGinLogger("account_delete"),
		a.deleteAccount,
	)

	// account auth, login
	authRouter := router.Group("/api/v1/account/auth")
	authRouter.POST("/",
		infra.SetGinLogger("account_auth_create_by_device"),
		a.loginByDevice,
	)
	authRouter.POST("/sso",
		infra.SetGinLogger("account_auth_create_by_sso"),
		a.loginBySSO,
	)
	authRouter.PUT("/",
		infra.SetGinLogger("account_auth_update_refresh_token"),
		a.refreshAuthToken,
	)
	authRouter.DELETE("/",
		infra.SetGinLogger("account_auth_delete"),
		a.logout,
	)

	return router
}

func (a *app) setDomainServices() {
	// db and cache connections
	dsn := config.GetContainerDSN()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("db.Ping on %s error: %w", dsn, err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetCacheAddress(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// repositories
	userRepo, err := user_repo.NewPostgresRepo(db)
	if err != nil {
		panic(fmt.Errorf("user_repo.NewPostgresRepo error: %w", err))
	}
	authRepo, err := auth_repo.NewRedisRepo(rdb)
	if err != nil {
		panic(fmt.Errorf("auth_repo.NewRedisRepo error: %w", err))
	}

	// services
	userSvc, err := user_svc.New(userRepo)
	if err != nil {
		panic(fmt.Errorf("user_svc.New error: %w", err))
	}
	authSvc := auth_svc.New(authRepo)

	a.userSvc = userSvc
	a.authSvc = authSvc
}

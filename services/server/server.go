package server

import (
	"context"
	"net/http"
	"os"
	"path"
	"service/services/bugsnag"
	"service/services/config"
	"service/services/health"
	"service/services/logger"
	"service/services/server/controllers"
	"service/services/shutdown"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var server http.Server

func init() {
	logger := logger.WithContext("server")
	config := config.QuietBuild(ServerConfig{})
	if config.AppEnv != "staging" && config.AppEnv != "production" {
		gin.SetMode(gin.DebugMode)
		router = gin.New()
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	}

	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/api/health", "/"),
		gin.Recovery(),
	)

	if config.AppEnv != "staging" && config.AppEnv != "production" {
		router.Use(cors.Default())
		cwd, _ := os.Getwd()
		router.Static("resources/html", path.Join(cwd, "resources/html"))
		logger.Info("Development mode activated, registering development routes...")
	} else {
		corsConfig := cors.Config{}
		allowOrigins := strings.Split(config.AllowOrigins, ",")
		allowOriginSet := make(map[string]bool)
		for _, v := range allowOrigins {
			allowOriginSet[v] = true
		}
		corsConfig.AllowOriginFunc = func(origin string) bool {
			_, originWhitelisted := allowOriginSet[origin]
			return originWhitelisted
		}
		logger.Info("Production mode activated, only allowing origins: %s", config.AllowOrigins)
		router.Use(cors.New(corsConfig))
	}

	// >>> all routes are registered here
	registerController(health.HealthController)

	server = http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}
}

func registerController[T any](controller controllers.Controller[T]) {
	var handlers []gin.HandlerFunc
	for _, middleware := range controller.Middlewares {
		handlers = append(handlers, middleware.AsGinMiddleware())
	}
	handlers = append(handlers, controller.Handler.AsGinHandler())

	switch controller.Method {
	case controllers.GET:
		router.GET(controller.Route, handlers...)
	case controllers.POST:
		router.POST(controller.Route, handlers...)
	case controllers.PUT:
		router.PUT(controller.Route, handlers...)
	case controllers.PATCH:
		router.PATCH(controller.Route, handlers...)
	case controllers.DELETE:
		router.DELETE(controller.Route, handlers...)
	}
}

var onFinishedShutdown func() = nil

func Start(shutdownHandler *shutdown.ShutdownHandler) {
	logger := logger.WithContext("server")

	onFinishedShutdown = shutdownHandler.NewShutdownHook("server", Stop)
	go func() {
		logger.Info("starting server at address %s", server.Addr)
		if startServerError := server.ListenAndServe(); startServerError != nil && startServerError != http.ErrServerClosed {
			logger.Error("failed to start server", bugsnag.FromError("Start Server Error", startServerError))
		}
	}()
}

func Stop() {
	logger := logger.WithContext("server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutdownError := server.Shutdown(ctx); shutdownError != nil {
		logger.Warn("shutdown error, ignoring: %s", shutdownError.Error())
	} else {
		logger.Info("shutdown finished")
	}

	if onFinishedShutdown != nil {
		onFinishedShutdown()
	}

}

package server

import (
	"context"
	"net/http"
	"os"
	"path"
	"service/services/bugsnag"
	"service/services/config"
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
	l := logger.For("server")
	c := config.QuietBuild(Config{})
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(
		logErrorRequests(),
		gin.Recovery(),
		controllers.AssignRequestId.AsGinMiddleware(),
		controllers.AssignRequestLanguage.AsGinMiddleware(),
	)

	if c.AppEnv != "staging" && c.AppEnv != "production" {
		router.Use(cors.Default())
		cwd, _ := os.Getwd()
		router.Static("resources/html", path.Join(cwd, "resources/html"))
		l.Info("Development mode activated, registering development routes...")
	} else {
		corsConfig := cors.Config{}
		allowOrigins := strings.Split(c.AllowOrigins, ",")
		allowOriginSet := make(map[string]bool)
		for _, v := range allowOrigins {
			allowOriginSet[v] = true
		}
		corsConfig.AllowOriginFunc = func(origin string) bool {
			_, originWhitelisted := allowOriginSet[origin]
			return originWhitelisted
		}
		l.Info("Production mode activated, only allowing origins: %s", c.AllowOrigins)
		router.Use(cors.New(corsConfig))
	}

	// >>> all routes are registered here
	registerRoutes()

	server = http.Server{
		Addr:    ":" + c.AppPort,
		Handler: router,
	}
}

func registerController[T any](controller controllers.Controller[T]) {
	l := logger.For("server")
	var handlers []gin.HandlerFunc
	for _, middleware := range controller.Middlewares {
		handlers = append(handlers, middleware.AsGinMiddleware())
	}
	handlers = append(handlers, controller.Handler.AsGinHandler())

	l.Info("%*s %*s --> %s", 5, controller.Method, 32, controller.Route, controller.HandlerTrace)
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

func Start(shutdownHandler *shutdown.Handler) {
	l := logger.For("server")

	onFinishedShutdown = shutdownHandler.NewShutdownHook("server", Stop)
	go func() {
		l.Info("starting server at address %s", server.Addr)
		if startServerError := server.ListenAndServe(); startServerError != nil && startServerError != http.ErrServerClosed {
			l.Error("failed to start server", bugsnag.FromError("Start Server Error", startServerError))
		}
	}()
}

func Stop() {
	l := logger.For("server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutdownError := server.Shutdown(ctx); shutdownError != nil {
		l.Warn("shutdown error, ignoring: %s", shutdownError.Error())
	} else {
		l.Info("shutdown finished")
	}

	if onFinishedShutdown != nil {
		onFinishedShutdown()
	}

}

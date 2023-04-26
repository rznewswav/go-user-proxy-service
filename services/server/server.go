package server

import (
	"context"
	"fmt"
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

func logErrorRequests() gin.HandlerFunc {
	logger := logger.WithContext("gin")

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		p := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		statusCode := c.Writer.Status()
		var loggerFn = logger.Debug
		var err error
		var contextError = c.Errors.ByType(gin.ErrorTypePrivate)
		if statusCode >= 500 {
			loggerFn = logger.Error
			err = fmt.Errorf("request returned server error")
		} else if len(contextError) > 0 {
			loggerFn = logger.Error
			err = fmt.Errorf("request returned context error: %s", contextError.String())
		}

		latency := time.Now().Sub(start)
		method := c.Request.Method
		url := p
		if len(raw) > 0 {
			url = fmt.Sprintf("%s?%s", p, raw)
		}

		if err != nil {
			loggerFn(
				"%s %s - %d (%dms)",
				method,
				url,
				statusCode,
				int64(latency.Seconds()),
				bugsnag.FromError("Request Error", err),
			)
		} else {
			loggerFn(
				"%s %s - %d (%dms)",
				method,
				url,
				statusCode,
				int64(latency.Seconds()),
			)
		}
	}
}

func init() {
	logger := logger.WithContext("server")
	config := config.QuietBuild(ServerConfig{})
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(
		logErrorRequests(),
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
	registerRoutes()

	server = http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}
}

func registerController[T any](controller controllers.Controller[T]) {
	l := logger.WithContext("server")
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

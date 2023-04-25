package server

import (
	"context"
	"net/http"
	"os"
	"path"
	"service/services/bugsnag"
	"service/services/config"
	"service/services/logger"
	server_routes "service/services/server/routes"
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
		router.GET("resources/images/:filename", server_routes.GetResourceImage)
		router.GET("resources/templates/:filename", server_routes.GetTemplateByName)
		router.POST("resources/templates/:filename", server_routes.PostSaveTemplateByName)
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
	router.GET("/api/health", server_routes.GetHealth)
	router.GET("/", server_routes.GetHealth)
	router.GET("/assets/:filename", server_routes.GetCovidImages)

	server = http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
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

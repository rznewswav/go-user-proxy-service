package redis

import (
	"context"
	"fmt"
	"service/services/config"
	"service/services/logger"
	"service/services/shutdown"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var initLock sync.Mutex
var initisLocked = false
var initExecError error
var onFinishedShutdown func()

func AddKeyPrefix(key string) string {
	return fmt.Sprintf("image-forge:%s", key)
}

func Connect(shutdownHandler *shutdown.ShutdownHandler) {
	if Client != nil {
		return
	}
	logger := logger.WithContext("redis")

	redisConfig := config.QuietBuild(RedisConfig{})

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.RedisHost, redisConfig.RedisPort),
		Password: redisConfig.RedisPassword,
	})

	onFinishedShutdown = shutdownHandler.NewShutdownHook("redis", QuietDisconnect)

	initLock.Lock()
	initisLocked = true

	go func() {
		defer func() {
			initLock.Unlock()
			initisLocked = false
		}()

		pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer pingCtxCancel()

		if pingError := Client.Ping(pingCtx).Err(); pingError != nil {
			logger.Error("redis ping returned error: %s", pingError)
			return
		}

		execCtx, execCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer execCtxCancel()
		if execError := Client.Get(execCtx, "random:key").Err(); execError != nil && execError != redis.Nil {
			logger.Error("redis get test returned error: %s", execError)
			initExecError = execError
			return
		}

	}()
}

func Health() error {
	if Client == nil {
		return ErrNotInitialized
	}

	if initExecError != nil {
		return ErrNotAvailable
	}

	pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCtxCancel()
	pingError := Client.Ping(pingCtx).Err()
	if pingError != nil {
		return fmt.Errorf("health check error: %w", pingError)
	}

	return nil
}

func Disconnect() error {
	logger := logger.WithContext("redis")
	if initisLocked {
		logger.Warn("redis is still trying to connect but a disconection was requested!")
	}

	initLock.Lock()
	logger.Debug("redis client lock requested and acquired, redis client is safe to be closed")
	initLock.Unlock()
	defer func() {
		if onFinishedShutdown != nil {
			logger.Debug("emitting shutdown finished hook")
			onFinishedShutdown()
		} else {
			logger.Debug("shutdown finished hook is not assigned")
		}
	}()

	if Client == nil {
		return nil
	}

	closeError := Client.Close()
	if closeError != nil {
		return closeError
	}

	Client = nil
	return nil
}

func QuietDisconnect() {
	logger := logger.WithContext("redis")
	disconnectError := Disconnect()
	if disconnectError != nil {
		logger.Error("error at disconnecting redis: %s", disconnectError.Error())
	}
}

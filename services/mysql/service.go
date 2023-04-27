package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service/services/config"
	"service/services/logger"
	"service/services/shutdown"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var initialized = false
var Client *sql.DB
var onFinishedShutdown func()
var initDbError error
var initExecError error
var initLock sync.Mutex
var initisLocked = false

func Connect(shutdownHandler *shutdown.Handler) {
	if initialized {
		return
	}
	l := logger.For("mysql")
	initialized = true
	dbConfig := config.QuietBuild(Config{})
	onFinishedShutdown = shutdownHandler.NewShutdownHook("mysql", QuietDisconnect)
	initLock.Lock()
	initisLocked = true
	go func() {
		defer func() {
			initLock.Unlock()
			initisLocked = false
		}()
		l.Info("entering goroutine...")
		if dbClient,
			dbConnectionError := sql.Open("mysql", dbConfig.DatabaseURL); dbConnectionError != nil {
			l.Error("mysql init returned error: %s", dbConnectionError)
			initDbError = dbConnectionError
			return
		} else {
			Client = dbClient
		}

		pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer pingCtxCancel()

		if pingError := Client.PingContext(pingCtx); pingError != nil {
			l.Error("mysql ping returned error: %s", pingError)
			return
		}

		execCtx, execCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer execCtxCancel()
		if _, execError := Client.ExecContext(execCtx, "SELECT id FROM corona_malaysia_weekly_state_stats LIMIT 1"); execError != nil {
			l.Error("mysql exec test returned error: %s", execError)
			initExecError = execError
			return
		}

		Client.SetConnMaxLifetime(3 * time.Minute)
		Client.SetMaxOpenConns(1)
		Client.SetMaxIdleConns(1)
	}()
}

var ErrNotInitialized = errors.New("client is not yet initialized")
var ErrNotAvailable = errors.New("client is not available, see logs")

func Health() error {
	if !initialized {
		return ErrNotInitialized
	}

	if initDbError != nil {
		return ErrNotAvailable
	}

	if initExecError != nil {
		return ErrNotAvailable
	}

	pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCtxCancel()
	pingError := Client.PingContext(pingCtx)
	if pingError != nil {
		return fmt.Errorf("health check error: %w", pingError)
	}

	return nil
}

func Disconnect() error {
	l := logger.For("mysql")
	if initisLocked {
		l.Warn("db is still trying to connect but a disconection was requested!")
	}

	initLock.Lock()
	l.Debug("db client lock requested and acquired, db client is safe to be closed")
	initLock.Unlock()
	defer func() {
		if onFinishedShutdown != nil {
			l.Debug("emitting shutdown finished hook")
			onFinishedShutdown()
		} else {
			l.Debug("shutdown finished hook is not assigned")
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
	l := logger.For("mysql")
	disconnectError := Disconnect()
	if disconnectError != nil {
		l.Error("error at disconnecting mysql: %s", disconnectError.Error())
	}
}

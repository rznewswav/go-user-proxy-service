package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"service/services/bugsnag"
	"service/services/logger"
	"time"
)

func logErrorRequests() gin.HandlerFunc {
	l := logger.For("gin")

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		p := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		statusCode := c.Writer.Status()
		var loggerFn = l.Debug
		var err error
		var contextError = c.Errors.ByType(gin.ErrorTypePrivate)
		if statusCode >= 500 {
			loggerFn = l.Error
			err = errors.New("request returned server error")
		} else if len(contextError) > 0 {
			loggerFn = l.Error
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

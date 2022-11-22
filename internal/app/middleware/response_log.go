package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"regexp"
	"time"
	"top-ping/pkg/logger"
	"top-ping/pkg/utils"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func ResponseLogger(config *logger.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		skipPaths := config.SkipPaths
		for _, skipPath := range skipPaths {
			reg := regexp.MustCompile(skipPath)
			if reg.MatchString(path) {
				c.Next()
				return
			}
		}

		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		cost := time.Since(start)
		responseBody := blw.body.String()
		if config.Desensitize {
			responseBody = utils.MaskJsonStr(&responseBody, config.SkipFields)
		}

		logger.Info(c.Request.Context(), "ResponseLog",
			zap.Int("Status", c.Writer.Status()),
			zap.String("Path", path),
			zap.String("Response", responseBody),
			zap.Duration("Cost", cost),
		)
	}
}

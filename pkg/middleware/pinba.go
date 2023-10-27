package middleware

import (
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/mkevac/gopinba"
	"time"
)

func PinbaHandler(config *structures.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.PinbaHost == "" {
			c.Next()

			return
		}

		start := time.Now()
		pc, err := gopinba.NewClient(config.PinbaHost)

		if err != nil {
			logger.LogError(err)
		}

		req := gopinba.Request{}

		req.Hostname = config.HostName
		req.ServerName = config.AppEnv
		req.ScriptName = c.Request.RequestURI
		req.Schema = "http"
		req.RequestCount = 1

		req.Tags = map[string]string{
			"type":   c.Request.Method,
			"method": c.Request.Method,
		}

		c.Next()

		req.Status = uint32(c.Writer.Status())
		req.RequestTime = time.Since(start)

		err = pc.SendRequest(&req)

		if err != nil {
			logger.LogError(err)
		}
	}
}

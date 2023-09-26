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
		req.ServerName = config.Name
		req.ScriptName = c.Request.RequestURI
		req.RequestCount = 1

		c.Next()

		req.Status = uint32(c.Writer.Status())
		req.RequestTime = time.Since(start)

		err = pc.SendRequest(&req)

		if err != nil {
			logger.LogError(err)
		}
	}
}

package middleware

import (
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ResponseHandler(config *structures.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			sentry.CaptureException(err)
			logError(err.Error(), c, config)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": err.Error()})

			return
		}

		appExceptionObject, exists := c.Get("exception")
		fmt.Printf("%+v\n", appExceptionObject)

		if !exists {
			data, _ := c.Get("data")
			logInfo("", c, config)
			c.JSON(http.StatusOK, gin.H{"success": true, "data": data})

			return
		}

		appException := exception.AppException{}
		mapstructure.Decode(appExceptionObject, &appException)
		sentry.CaptureException(appException.Error)
		fmt.Printf("%+v\n", appException)
		logError(appException.Error.Error(), c, config)
		c.JSON(appException.Code, gin.H{"success": false, "message": appException.Error.Error(), "details": appException.Context})
	}
}

func logInfo(message string, c *gin.Context, config *structures.AppConfig) {
	logResponse("INFO", message, c, config)
}

func logError(message string, c *gin.Context, config *structures.AppConfig) {
	logResponse("ERROR", message, c, config)
}

func logResponse(level string, message string, c *gin.Context, config *structures.AppConfig) {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	logData := make([]string, 5)
	logData[0] = time.Now().Format("2006-01-02 15:04:05.345")
	logData[1] = level
	logData[2] = "[" + c.Request.Method + ", " + c.Request.RequestURI + ", " + strconv.Itoa(c.Writer.Status()) + "]"
	logData[3] = "[" + config.Name + ", " + c.GetHeader("X-B3-TraceId") + "]"
	logData[4] = message

	log.Println(strings.Join(logData, " "))
}

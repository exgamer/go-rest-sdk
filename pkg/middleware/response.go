package middleware

import (
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func ResponseHandler(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": err.Error()})

		return
	}

	appExceptionObject, exists := c.Get("exception")
	fmt.Printf("%+v\n", appExceptionObject)

	if !exists {
		data, _ := c.Get("data")
		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})

		return
	}

	appException := exception.AppException{}
	mapstructure.Decode(appExceptionObject, &appException)
	fmt.Printf("%+v\n", appException)
	c.JSON(appException.Code, gin.H{"success": false, "message": appException.Error.Error()})
}

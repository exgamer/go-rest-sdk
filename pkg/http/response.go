package http

import (
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, exception *exception.AppException) {
	c.Set("exception", exception)
}

func Response(c *gin.Context, status int, data any) {
	c.Set("data", data)
}

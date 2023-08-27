package gin

import (
	"github.com/exgamer/go-rest-sdk/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// Options
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middleware.ResponseHandler) // обработчик ошибок
	router.Use(gin.Recovery())
	router.Use(middleware.Recovery())

	return router
}

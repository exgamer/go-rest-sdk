package gin

import (
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
)

func InitRouter() *gin.Engine {
	// Options
	router := gin.Default()
	router.Use(gin.Logger())
	//router.Use(gin.Recovery())
	//router.Use(middleware.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "405 method not allowed"})
	})

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	return router
}

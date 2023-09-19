package gin

import (
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
)

func InitRouter(appConfig *structures.AppConfig) *gin.Engine {
	// Options
	router := gin.Default()
	router.Use(gin.Logger())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "405 method not allowed"})
	})

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: appConfig.SentryDsn,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	router.Use(sentrygin.New(sentrygin.Options{}))

	return router
}

package api

import (
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/config"
	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf       *config.Config
	Middleware *middleware.Middleware
	Handler    *handlers.Handler
}

func NewRouter(opt *Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), opt.Middleware.Logging())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return router
}

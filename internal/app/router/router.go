package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"top-ping/internal/app/middleware"
	"top-ping/pkg/logger"
	"top-ping/pkg/rest"
	"top-ping/pkg/utils"
)

func Router(profile string, logging *logger.Config) *gin.Engine {
	if profile == utils.ProdProfile {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	var r = gin.New()
	if profile != utils.ProdProfile {
		pprof.Register(r)
		//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Use(middleware.AccessLogger(logging))
	r.Use(middleware.ResponseLogger(logging))
	r.Use(gin.Recovery())

	apiV1 := r.Group("/v1")
	apiV1.Use()

	{
		//apiV1.POST("/user/get_one", controller.GetUser)
		//apiV1.POST("/student/get_one", controller.GetStudent)
	}

	r.GET("/ping", func(c *gin.Context) {
		logger.Info(c.Request.Context(), "ping")
		rest.R.Success(c, "pong")
	})

	return r
}

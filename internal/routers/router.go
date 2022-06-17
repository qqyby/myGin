package routers

import (
	"github.com/gin-gonic/gin"
	"myGin/global"
	"myGin/internal/routers/middlewares"
	"myGin/settings"
	"net/http"
)

func Setup() *gin.Engine  {
	r := gin.New()

	r.Use(middlewares.RequestId())
	if settings.AppCfg.RunMode == global.DebugModel {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		gin.SetMode(gin.ReleaseMode)
		r.Use(middlewares.Logger())
		r.Use(middlewares.Recovery(true))
	}

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	v1 := r.Group("/group")

	v1.POST("")
	v1.GET("")
	v1.PUT("")
	v1.DELETE("")

	return r
}
package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine  {
	r := gin.New()

	r.Use()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  404,
		})
	})

	v1 := r.Group("/group")

	v1.POST("")
	v1.GET("")
	v1.PUT("")
	v1.DELETE("")

	return r
}
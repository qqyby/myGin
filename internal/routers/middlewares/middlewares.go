package middlewares

import "github.com/gin-gonic/gin"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.SetRequestId(c)
		c.Next()
	}
}

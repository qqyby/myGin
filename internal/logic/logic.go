package logic

import (
	"github.com/gin-gonic/gin"
	"myGin/pkg/app"
)

type Logic struct {
	Ctx       *gin.Context
	RequestId string
	SecretId  string
	UserId    string
}

func New(c *gin.Context) *Logic {
	return &Logic{
		Ctx:       c,
		RequestId: app.GetRequestId(c),
	}
}

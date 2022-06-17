package app

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const (
	KeyRequestID = "RequestId" // RequestId 对应的键名
	KeySecretID  = "SecretId"  // 密钥对id 对应的键名
	KeyUseID     = "UseId"
)

func SetRequestId(c *gin.Context) string {
	requestId := uuid.NewV4().String()
	c.Set(KeyRequestID, requestId)
	return requestId
}

func GetRequestId(c *gin.Context) string {
	requestId, ok := c.Get(KeyRequestID)
	if !ok {
		return ""
	}
	return requestId.(string)
}

func SetSecretId(c *gin.Context, secretId string) string {
	c.Set(KeySecretID, secretId)
	return secretId
}

func GetSecretId(c *gin.Context) string {
	secretId, ok := c.Get(KeySecretID)
	if !ok {
		return ""
	}
	return secretId.(string)
}

func SetUserId(c *gin.Context, uId string) string {
	c.Set(KeyUseID, uId)
	return uId
}

func GetUserId(c *gin.Context) string {
	useId, ok := c.Get(KeyUseID)
	if !ok {
		return ""
	}
	return useId.(string)
}

func ZapRequestId(c *gin.Context) zap.Field {
	requestId := GetRequestId(c)
	return zap.String(KeyRequestID, requestId)
}

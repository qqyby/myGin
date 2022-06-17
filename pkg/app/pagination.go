package app

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func TrimQuery(c *gin.Context, key string) string {
	return strings.Trim(c.Query(key), " ")
}

func GetPage(c *gin.Context) int {
	page := StrTo(TrimQuery(c, "page")).MustInt()
	if page <= 0 {
		return 1
	}
	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := StrTo(TrimQuery(c, "page_size")).MustInt()
	if pageSize <= 0 {
		return DefaultPageSize
	}

	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	var result int
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}

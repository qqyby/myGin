package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"myGin/pkg/app"
	"myGin/pkg/errcode"
	"myGin/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.SetRequestId(c)
		c.Next()
	}
}

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}


// Logger 接收gin框架默认的日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		requestId := app.SetRequestId(c)
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.String("request_id", requestId),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}


func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := app.NewResponse(c)
		urlPath := strings.TrimRight(c.Request.URL.Path, "/")

		if err := checkLicense(urlPath); err != nil {
			zap.L().Error("jwt: check license failed", zap.String("url", urlPath), app.PlainError(err))
			response.ToErrorResponse(errcode.LicenseExpiredError)
			c.Abort()
			return
		}

		// 不开启认证，后端记录使用admin
		if !settings.JwtCfg.Enable {
			// 将默认超级管理员存储到上下文中
			logic.SaveJwtToGinContext(c, logic.NewJwTAdminClaims())
			c.Next()
			return
		}

		var token string
		if t := app.TrimQuery(c, "token"); t != "" {
			token = t
		} else {
			token = c.GetHeader("token")
		}

		if token == "" {
			zap.L().Error("jwt: token empty", zap.String("url", urlPath))
			response.ToErrorResponse(errcode.UnauthorizedTokenError)
			c.Abort()
			return
		}

		Claims, err := logic.ParseToken(token)
		if err != nil {
			zap.L().Error("jwt: Parse Token error", zap.String("url", urlPath), app.PlainError(err))
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				response.ToErrorResponse(errcode.UnauthorizedTokenTimeout)
				c.Abort()
				return
			default:
				response.ToErrorResponse(errcode.UnauthorizedTokenError)
				c.Abort()
				return
			}
		}

		// 存储到上下文中
		logic.SaveJwtToGinContext(c, Claims)
		//获取菜单,修改自己的密码，分类管理 不需要进一步验证权限
		if urlPath == "/api/v1/menu_list" || urlPath == "/api/v1/auth/users/my_pwd" || strings.HasPrefix(urlPath, "/api/v1/categories") {
			c.Next()
			return
		}

		// 超级管理员
		if Claims.RoleRow.IsAdmin {
			c.Next()
			return
		}

		// 是否有权限
		var authorized bool
		for _, p := range Claims.RoleRow.Privileges {
			if strings.HasPrefix(urlPath, p.MenuPath) {
				authorized = true
				break
			}
		}

		if !authorized {
			zap.L().Error("jwt: request is unauthorized", zap.String("url", urlPath))
			response.ToErrorResponse(errcode.UnauthorizedToAccess)
			c.Abort()
			return
		}
		c.Next()
	}
}


func checkLicense(urlPath string) error {
	if settings.LicenseCfg.Expire() {
		return fmt.Errorf("jwt: request: %v; license expired", urlPath)
	}

	//if settings.LicenseCfg.Version == global.LicenseCustomerVersion && urlPath != "/api/v1/menu_list" && urlPath != "/api/v1/sys/license" &&
	//	urlPath != "/api/v1/sys/hardware" && urlPath != "/api/v1/sys/restart_server" {
	//	return fmt.Errorf("jwt: request: %v; license just for customer, can not do this", urlPath)
	//}
	return nil
}
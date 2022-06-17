package logic

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"myGin/internal/model"
	"myGin/settings"
	"time"
)

const (
	JWTSecret = "eyJhbGciOiJIUzI1NiIsJIUzI1NiIsInR5cCI6IkpXVCJ9"
	JWTIssuer = "bravo"
	KeyJWT    = "JwtClaims" // jwt在gin 上下文对应的键名
)

type RoleRow struct {
	RoleID     string                `json:"role_id"`
	RoleName   string                `json:"role_name"`
	IsAdmin    bool                  `json:"is_admin"`
	Remark     string                `json:"remark"`
	CreateTime string                `json:"create_time"`
	Privileges []*model.PrivilegeRow `json:"privileges"`
}

type Claims struct {
	UserID   string   `json:"user_id"`
	UserName string   `json:"user_name"`
	RoleRow  *RoleRow `json:"role_row"`

	// jwt-go 预定义的内容
	jwt.StandardClaims
}

func (c *Claims) UserInfo() string {
	return fmt.Sprintf("%s-%s", c.UserID, c.UserName)
}

// 获取该项目的 jwt Secret
func GetJWTSecret() []byte {
	return []byte(JWTSecret)
}

func GenerateToken(userID, userName string, roleRow *RoleRow) (string, error) {
	expireTime := time.Now().Add(settings.JwtCfg.Expire * time.Hour)
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RoleRow:  roleRow,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    JWTIssuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// 不开启验证的情况下，默认使用超级管理源，admin来进行初始化
func NewJwTAdminClaims() *Claims {
	claims := Claims{
		UserID:   "5161828556476416",
		UserName: "admin",
		RoleRow: &RoleRow{
			RoleID:   "5161611304112128",
			RoleName: "超级管理员",
			IsAdmin:  true,
		},
	}
	return &claims
}

func SaveJwtToGinContext(c *gin.Context, claims *Claims) {
	c.Set(KeyJWT, claims)
}

func GetJwtFromGinContext(c *gin.Context) (*Claims, error) {
	value, ok := c.Get(KeyJWT)
	if !ok {
		return nil, fmt.Errorf("get jwt from context failed")
	}

	claims, ok := value.(*Claims)
	if !ok {
		return nil, fmt.Errorf("transfer jwt failed")
	}

	return claims, nil
}

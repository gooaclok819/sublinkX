package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"sublink/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 随机密钥

// var Secret = []byte("sublink") // 秘钥
var Secret = []byte(models.ReadConfig().JwtSecret) // 从配置文件读取JWT密钥

// JwtClaims jwt声明
type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AuthorToken 验证token中间件
func AuthorToken(c *gin.Context) {
	// 定义白名单
	list := []string{"/static", "/api/v1/auth/login", "/api/v1/auth/captcha", "/c/", "/api/v1/version"}
	// 如果是首页直接跳过
	if c.Request.URL.Path == "/" {
		c.Next()
		return
	}
	// 如果是白名单直接跳过
	for _, v := range list {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			c.Next()
			return
		}
	}
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{"msg": "请求未携带token"})
		c.Abort()
		return
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		c.JSON(400, gin.H{"msg": "token格式错误"})
		c.Abort()
		return
	}
	// 去掉Bearer前缀
	token = strings.Replace(token, "Bearer ", "", -1)
	mc, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		c.Abort()
		return
	}
	c.Set("username", mc.Username)
	c.Next()
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JwtClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

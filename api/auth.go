package api

import (
	"log"
	"sublink/middlewares"
	"sublink/models"
	"sublink/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 获取token
func GetToken(username string) (string, error) {
	// 过期时间天
	ExpireDays := models.ReadConfig().ExpireDays
	c := &middlewares.JwtClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(ExpireDays)).Unix(), // 设置过期时间
			IssuedAt:  time.Now().Unix(),                                                 // 签发时间
			Subject:   username,                                                          // 用户
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(middlewares.Secret)
}

// 获取captcha图形验证码
func GetCaptcha(c *gin.Context) {
	id, bs4, _, err := utils.GetCaptcha()
	if err != nil {
		log.Println("获取验证码失败")
		c.JSON(400, gin.H{
			"msg": "获取验证码失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"captchaKey":    id,
			"captchaBase64": bs4,
		},
		"msg": "获取验证码成功",
	})

}

// 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captchaCode := c.PostForm("captchaCode")
	captchaKey := c.PostForm("captchaKey")
	// 验证验证码
	if !utils.VerifyCaptcha(captchaKey, captchaCode) {
		log.Println("验证码错误")
		c.JSON(400, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	user := &models.User{Username: username, Password: password}
	err := user.Verify()
	if err != nil {
		log.Println("账号或者密码错误")
		c.JSON(400, gin.H{
			"msg": "账号或者密码错误",
		})
		return
	}
	// 生成token
	token, err := GetToken(username)
	if err != nil {
		log.Println("获取token失败", err)
		c.JSON(400, gin.H{
			"msg": "获取token失败",
		})
		return
	}
	// 登录成功返回token
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"accessToken":  token,
			"tokenType":    "Bearer",
			"refreshToken": nil,
			"expires":      nil,
		},
		"msg": "登录成功",
	})
}
func UserOut(c *gin.Context) {
	// 拿到jwt中的username
	if _, Is := c.Get("username"); Is {
		c.JSON(200, gin.H{
			"code": "00000",
			"msg":  "退出成功",
		})
	}
}

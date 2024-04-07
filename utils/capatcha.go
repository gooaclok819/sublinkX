package utils

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GetCaptcha 获取验证码
func GetCaptcha() (string, string, string, error) {
	driver := base64Captcha.NewDriverMath(60, 180, 80, 0, &color.RGBA{255, 255, 255, 255}, nil, nil)
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, answer string) bool {
	return store.Verify(id, answer, true)
}

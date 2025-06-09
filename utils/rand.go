package utils

import (
	"math/rand"
)

// RandString 生成随机字符串
func RandString(number int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	// 用 []byte 直接构造字符串
	n := rand.Intn(number) + 1 // 防止生成空字符串，范围是1到31
	randomString := make([]byte, n)
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(str))
		randomString[i] = str[randomIndex]
	}
	Secret := string(randomString)
	return Secret
}

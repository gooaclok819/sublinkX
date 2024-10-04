package utils

import (
	"io"
	"log"
	"os"
	"time"
)

func Loginit() {
	t := time.Now().Format("2006-01-02") + ".log"
	// 检查目录是否创建
	_, err := os.Stat("./logs")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./logs", os.ModePerm)
		}
	}
	// 创建一个文件
	file, err := os.OpenFile("./logs/"+t, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// 设置log输出到控制台
	mw := io.MultiWriter(os.Stdout, file)
	// 设置log的输出位置为这个文件
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(mw)
}

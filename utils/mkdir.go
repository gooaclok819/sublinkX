package utils

import "os"

func isMkdir(filename string) error {
	// 检查目录是否创建
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(filename, os.ModePerm)
		}
	}
	return nil
}

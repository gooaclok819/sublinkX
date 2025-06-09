package models

import (
	"fmt"
	"log"
	"os"
	"sublink/utils"

	"gopkg.in/yaml.v3"
)

// type Config struct {
// 	ID    int
// 	Key   string
// 	Value string
// }

// Config 配置结构体
type Config struct {
	JwtSecret  string `yaml:"jwt_secret"`  // JWT密钥
	ExpireDays int    `yaml:"expire_days"` // 过期天数
	Port       int    `yaml:"port"`        // 端口号
}

var comment string = `# jwt_secret: JWT密钥
# expire_days: token 过期天数
# port: 启动端口
`

// 初始化配置
func ConfigInit() {
	// 检查配置文件是否存在
	if _, err := os.Stat("./db/config.yaml"); os.IsNotExist(err) {
		R := utils.RandString(31) // 生成随机字符串作为JWT密钥
		// 如果不存在则创建默认配置文件
		defaultConfig := Config{
			JwtSecret:  R, // 生成随机JWT密钥
			ExpireDays: 14,
			Port:       8000, // 默认端口
		}

		// 生成yaml文件
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			log.Println("生成默认配置文件失败:", err)
			return
		}
		data = []byte(comment + string(data)) // 添加注释
		err = os.WriteFile("./db/config.yaml", data, 0644)
		if err != nil {
			fmt.Println("写入文件失败:", err)
			return
		}
		log.Println("配置文件不存在，已创建默认配置文件")
	}
}

// 读取配置
func ReadConfig() Config {
	file, err := os.ReadFile("./db/config.yaml")
	if err != nil {
		log.Println(err)
	}
	cfg := Config{}
	yaml.Unmarshal(file, &cfg)
	return cfg
}

// 设置配置
func SetConfig(newCfg Config) {
	oldCfg := ReadConfig() // 读取旧的配置文件
	// 覆盖新的字段
	if newCfg.JwtSecret != "" {
		oldCfg.JwtSecret = newCfg.JwtSecret
	}
	if newCfg.ExpireDays != 0 {
		oldCfg.ExpireDays = newCfg.ExpireDays
	}
	if newCfg.Port != 0 {
		oldCfg.Port = newCfg.Port
	}
	// 写入文件
	data, err := yaml.Marshal(&oldCfg)
	if err != nil {
		log.Println(err)
	}
	data = []byte(comment + string(data)) // 添加注释
	os.WriteFile("./db/config.yaml", data, 0644)
}

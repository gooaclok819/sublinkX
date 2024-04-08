package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sublink/middlewares"
	"sublink/models"
	"sublink/routers"
	"sublink/utils"

	"github.com/gin-gonic/gin"
)

//go:embed static/js/*
//go:embed static/css/*
//go:embed static/img/*
//go:embed static/index.html
var embeddedFiles embed.FS

//go:embed template/clash.yaml
var clashTemplate []byte

func Templateinit() {
	// 设置template路径

	// 检查目录是否创建
	_, err := os.Stat("./template")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./template", os.ModePerm)
		}
	}
	_, err = os.Stat("./template/clash.yaml")
	if os.IsNotExist(err) {
		err = os.WriteFile("./template/clash.yaml", clashTemplate, 0666)
		if err != nil {
			log.Println(err)
			return
		}
	}

}
func main() {
	// 初始化gin框架
	r := gin.Default()
	// 初始化日志配置
	utils.Loginit()
	// 初始化数据库
	models.InitSqlite()
	// 初始化模板
	Templateinit()
	// 安装中间件
	r.Use(middlewares.AuthorToken) // jwt验证token
	// 设置静态资源路径
	staticFiles, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Println(err)
	}
	r.StaticFS("/static", http.FS(staticFiles))
	// 设置模板路径
	r.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(staticFiles, "index.html")
		if err != nil {
			c.Error(err)
			return

		}
		c.Data(200, "text/html", data)
	})
	// 注册路由
	routers.User(r)
	routers.Mentus(r)
	routers.Subcription(r)
	routers.Nodes(r)
	routers.Clients(r)
	routers.Total(r)
	// 启动服务
	r.Run(":8000")
}

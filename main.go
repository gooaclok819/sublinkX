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

//go:embed template
var Template embed.FS

func Templateinit() {
	// 设置template路径

	// 检查目录是否创建
	Template, err := fs.Sub(embeddedFiles, "template")
	if err != nil {
		log.Println(err)
	}
	Templates, err := fs.ReadDir(Template, ".")
	if err != nil {
		log.Println(err)
	}
	for _, v := range Templates {
		_, err := os.Stat("./template/" + v.Name())
		//如果文件不存在则写入默认模板
		if os.IsNotExist(err) {
			data, err := fs.ReadFile(Template, v.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			err = os.WriteFile("./template/"+v.Name(), data, 0666)
			if err != nil {
				log.Println(err)
			}
		}
	}
	// _, err := os.Stat("./template/clash.yaml")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		os.MkdirAll("./template", os.ModePerm)
	// 	}
	// }
	// _, err = os.Stat("./template/clash.yaml")
	// if os.IsNotExist(err) {
	// 	err = os.WriteFile("./template/clash.yaml", clashTemplate, 0666)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }

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
	routers.Templates(r)
	// 启动服务
	r.Run(":8000")
}

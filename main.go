package main

import (
	"sublink/middlewares"
	"sublink/models"
	"sublink/routers"
	"sublink/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化gin框架
	r := gin.Default()
	// 初始化日志配置
	utils.Loginit()
	// 初始化数据库
	models.InitSqlite()
	// 安装中间件
	r.Use(middlewares.AuthorToken) // jwt验证token
	// 设置静态资源路径
	r.Static("/static", "./static")
	// 设置模板路径
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	// 注册路由
	routers.User(r)
	routers.Mentus(r)
	routers.Subcription(r)
	routers.Nodes(r)
	routers.Clients(r)
	routers.Total(r)
	// 启动服务
	r.Run(":8080")
}

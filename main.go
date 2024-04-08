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
	// 注册路由
	routers.User(r)
	routers.Mentus(r)
	routers.Subcription(r)
	routers.Nodes(r)
	routers.Clients(r)
	// 启动服务
	r.Run(":8080")
}

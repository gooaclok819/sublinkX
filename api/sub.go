// api/subcription.go

package api

import (
	// 导入 json 包，用于解析 config 字符串

	"log"
	"strconv"
	"strings"
	"sublink/models" // 导入 models 包

	"github.com/gin-gonic/gin"
)

func SubTotal(c *gin.Context) {
	var Sub models.Subcription
	subs, err := Sub.List()
	count := len(subs)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "取得订阅总数失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": count,
		"msg":  "取得订阅总数",
	})
}

// 获取订阅列表
func SubGet(c *gin.Context) {
	var Sub models.Subcription
	Subs, err := Sub.List()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "node list error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": Subs,
		"msg":  "node get",
	})
}

// 添加订阅
func SubAdd(c *gin.Context) {
	name := c.PostForm("name")
	configs := c.PostForm("config") // 这里的 configString 是前端传来的 JSON 字符串
	nodes := c.PostForm("nodes")

	if name == "" || nodes == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称或节点不能为空",
		})
		return
	}

	// 1. 根据 nodesString 字符串，构建 models.Node 数组
	var NodesData []models.Node

	for _, nodeName := range strings.Split(nodes, ",") {
		if strings.TrimSpace(nodeName) == "" {
			continue
		}
		FirstNode := models.Node{
			Name: nodeName,
		}

		// 查出node的数据
		result := models.DB.Model(models.Node{}).Where("name = ?", FirstNode.Name).First(&FirstNode)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(400, gin.H{
				"msg": result.Error,
			})
			return
		}
		// 插入nodes
		NodesData = append(NodesData, FirstNode)
	}
	sub := models.Subcription{
		Name:      name,
		Config:    configs,   // 这里直接赋值字符串
		NodeOrder: nodes,     // 这里直接赋值字符串
		Nodes:     NodesData, // 这里直接赋值 nodes 数组

	}
	err := sub.Add()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "添加订阅失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "添加订阅成功",
	})
}

// 更新订阅
func SubUpdate(c *gin.Context) {
	NewName := c.PostForm("name")
	OldName := c.PostForm("oldname")
	configs := c.PostForm("config") // 这里的 configString 是前端传来的 JSON 字符串
	nodes := c.PostForm("nodes")

	if NewName == "" || nodes == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称或节点不能为空",
		})
		return
	}

	// 1. 根据 nodesString 字符串，构建 models.Node 数组
	var NodesData []models.Node

	for _, nodeName := range strings.Split(nodes, ",") {
		if strings.TrimSpace(nodeName) == "" {
			continue
		}
		FirstNode := models.Node{
			Name: nodeName,
		}

		// 查出node的数据
		result := models.DB.Model(models.Node{}).Where("name = ?", FirstNode.Name).First(&FirstNode)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(400, gin.H{
				"msg": result.Error,
			})
			return
		}
		// 插入nodes
		NodesData = append(NodesData, FirstNode)
	}
	OldSub := models.Subcription{
		Name: OldName,
	}
	NewSub := models.Subcription{
		Name:      NewName,
		Config:    configs,   // 这里直接赋值字符串
		NodeOrder: nodes,     // 这里直接赋值字符串
		Nodes:     NodesData, // 这里直接赋值 nodes 数组

	}

	err := OldSub.Update(&NewSub)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "更新订阅失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "更新订阅成功",
	})
}

// 删除订阅 (无需修改)
func SubDel(c *gin.Context) {
	var sub models.Subcription
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, err := strconv.Atoi(id) // 增加错误检查
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "无效的 ID: " + err.Error(),
		})
		return
	}
	sub.ID = x
	err = sub.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "查找订阅失败: " + err.Error(),
		})
		return
	}
	err = sub.Del()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "删除订阅失败: " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "删除订阅成功",
	})
}

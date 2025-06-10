// api/subcription.go

package api

import (
	// 导入 json 包，用于解析 config 字符串
	"log"
	"strconv"
	"strings"
	"sublink/models" // 导入 models 包
	"time"

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
	var sub models.Subcription
	name := c.PostForm("name")
	configString := c.PostForm("config") // 这里的 configString 是前端传来的 JSON 字符串
	nodesString := c.PostForm("nodes")

	if name == "" || nodesString == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称或节点不能为空",
		})
		return
	}

	// 1. 根据 nodesString 字符串，构建 models.Node 数组
	var selectedNodes []models.Node
	for _, nodeName := range strings.Split(nodesString, ",") {
		trimmedName := strings.TrimSpace(nodeName)
		if trimmedName == "" {
			continue
		}
		var node models.Node
		node.Name = trimmedName
		err := node.Find()
		if err != nil {
			log.Printf("Warning: Node with name '%s' not found for subscription '%s'. Skipping.", trimmedName, name)
			continue
		}
		selectedNodes = append(selectedNodes, node)
	}
	sub.Nodes = selectedNodes

	// 2. 将前端传来的原始排序字符串赋值给 NodeOrder 字段
	sub.NodeOrder = nodesString

	// 3. 将前端传来的 config JSON 字符串直接赋值给 sub.Config
	sub.Config = configString // <--- 这里直接赋值字符串

	sub.Name = name
	sub.CreateDate = time.Now().Format("2006-01-02 15:04:05")
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
	var sub models.Subcription
	name := c.PostForm("name")
	oldname := c.PostForm("oldname")
	configString := c.PostForm("config") // 这里的 configString 是前端传来的 JSON 字符串
	nodesString := c.PostForm("nodes")

	if name == "" || nodesString == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称或节点不能为空",
		})
		return
	}

	// 1. 查找要更新的旧订阅
	sub.Name = oldname
	err := sub.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "查找旧订阅失败: " + err.Error(),
		})
		return
	}

	// 2. 根据 nodesString 字符串，构建新的 models.Node 数组
	var newSelectedNodes []models.Node
	for _, nodeName := range strings.Split(nodesString, ",") {
		trimmedName := strings.TrimSpace(nodeName)
		if trimmedName == "" {
			continue
		}
		var node models.Node
		node.Name = trimmedName
		err := node.Find()
		if err != nil {
			log.Printf("Warning: Node with name '%s' not found for subscription '%s'. Skipping.", trimmedName, name)
			continue
		}
		newSelectedNodes = append(newSelectedNodes, node)
	}

	// 3. 更新 Subcription 字段
	sub.Config = configString // <--- 这里直接赋值字符串
	sub.Name = name
	sub.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	sub.Nodes = newSelectedNodes
	sub.NodeOrder = nodesString

	// 4. 调用 sub.Update() 方法
	err = sub.Update()
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

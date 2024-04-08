package api

import (
	"strconv"
	"strings"
	"sublink/models"
	"time"

	"github.com/gin-gonic/gin"
)

func NodeUpdadte(c *gin.Context) {
	var node models.Node
	name := c.PostForm("name")
	link := c.PostForm("link")
	if name == "" || link == "" {
		c.JSON(400, gin.H{
			"msg": "节点名称 or 备注不能为空",
		})
		return
	}
	node.Name = name
	node.Link = link
	err := node.Update()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "更新失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "更新成功",
	})
}

// 获取节点列表
func NodeGet(c *gin.Context) {
	var Node models.Node
	nodes, err := Node.List()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "node list error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": nodes,
		"msg":  "node get",
	})
}

// 添加节点
func NodeAdd(c *gin.Context) {
	var Node models.Node
	link := c.PostForm("link")
	name := c.PostForm("name")
	if link == "" || name == "" {
		c.JSON(400, gin.H{
			"msg": "link or name 不能为空",
		})
		return
	}
	if !strings.Contains(link, "://") {
		c.JSON(400, gin.H{
			"msg": "link 必须包含 ://",
		})
		return
	}
	Node.Link = link
	Node.Name = name
	Node.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	err := Node.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = Node.Add()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "添加失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "添加成功",
	})
}

// 删除节点
func NodeDel(c *gin.Context) {
	var Node models.Node
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, _ := strconv.Atoi(id)
	Node.ID = x
	err := Node.Del()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "删除失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "删除成功",
	})
}

// 节点统计
func NodesTotal(c *gin.Context) {
	var Node models.Node
	nodes, err := Node.List()
	count := len(nodes)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "获取不到节点统计",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": count,
		"msg":  "取得节点统计",
	})
}

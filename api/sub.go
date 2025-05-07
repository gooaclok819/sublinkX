package api

import (
	"github.com/goccy/go-json"
	"strconv"
	"strings"
	"sublink/models"
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

// 添加节点
func SubAdd(c *gin.Context) {
	var sub models.Subcription
	name := c.PostForm("name")
	config := c.PostForm("config")
	nodes := c.PostForm("nodes")
	if name == "" || nodes == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称 or 节点不能为空",
		})
		return
	}
	sub.Nodes = []models.Node{}
	for _, v := range strings.Split(nodes, ",") {
		var node models.Node
		node.Name = v
		err := node.Find()
		if err != nil {
			continue
		}
		sub.Nodes = append(sub.Nodes, node)
	}

	sub.Config = config
	sub.Name = name
	sub.CreateDate = time.Now().Format("2006-01-02 15:04:05")

	err := sub.Add()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "添加失败",
		})
		return
	}
	err = sub.AddNode() //创建多对多关系
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "添加成功",
	})
}

// 更新节点
func SubUpdate(c *gin.Context) {
	var sub models.Subcription
	name := c.PostForm("name")
	oldname := c.PostForm("oldname")
	config := c.PostForm("config")
	nodes := c.PostForm("nodes")
	if name == "" || nodes == "" {
		c.JSON(400, gin.H{
			"msg": "订阅名称 or 节点不能为空",
		})
		return
	}
	// 查找旧节点
	sub.Name = oldname
	err := sub.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 更新节点
	sub.Config = config
	sub.Name = name
	sub.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	sub.Nodes = []models.Node{}
	for _, v := range strings.Split(nodes, ",") {
		var node models.Node
		node.Name = v
		err := node.Find()
		if err != nil {
			continue
		}
		sub.Nodes = append(sub.Nodes, node)
	}

	err = sub.Update()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "更新失败",
		})
		return
	}

	err = sub.UpdateNodes() //更新多对多关系
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "更新成功",
	})
}

// 删除节点
func SubDel(c *gin.Context) {
	var sub models.Subcription
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, _ := strconv.Atoi(id)
	sub.ID = x
	err := sub.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "查找失败",
		})
		return
	}
	err = sub.Del()
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

// 更新节点排序
func SubUpdateSort(c *gin.Context) {
	var nodeSorts []models.SubscriptionNodes

	subIDStr := c.PostForm("subId")
	if subIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "订阅ID不能为空",
		})
		return
	}

	subID, err := strconv.Atoi(subIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "无效的订阅ID",
		})
		return
	}

	// 获取节点排序数据
	nodesData := c.PostForm("nodes")
	if nodesData == "" {
		c.JSON(400, gin.H{
			"msg": "节点数据不能为空",
		})
		return
	}

	// 解析节点排序数据
	if err := json.Unmarshal([]byte(nodesData), &nodeSorts); err != nil {
		c.JSON(400, gin.H{
			"msg": "解析节点数据失败",
		})
		return
	}
	sn := &models.SubscriptionNodes{}
	if err := sn.UpdateNodesWithSort(subID, nodeSorts); err != nil {
		c.JSON(500, gin.H{
			"msg": "更新排序失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "排序更新成功",
	})
}

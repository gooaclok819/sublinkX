package api

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sublink/models"
	"sublink/node"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DocodeNodeName(nd *models.Node) (models.Node, error) { // 解码节点名称
	if nd.Name == "" {
		u, err := url.Parse(nd.Link)
		if err != nil {
			log.Println(err)
			return *nd, err
		}
		switch {
		case u.Scheme == "ss":
			ss, err := node.DecodeSSURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = ss.Name
		case u.Scheme == "ssr":
			ssr, err := node.DecodeSSRURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = ssr.Qurey.Remarks
		case u.Scheme == "trojan":
			trojan, err := node.DecodeTrojanURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = trojan.Name
		case u.Scheme == "vmess":
			vmess, err := node.DecodeVMESSURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = vmess.Ps

		case u.Scheme == "vless":
			vless, err := node.DecodeVLESSURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = vless.Name
		case u.Scheme == "hy" || u.Scheme == "hysteria":
			hy, err := node.DecodeHYURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = hy.Name
		case u.Scheme == "hy2" || u.Scheme == "hysteria2":
			hy2, err := node.DecodeHY2URL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = hy2.Name
		case u.Scheme == "tuic":
			tuic, err := node.DecodeTuicURL(nd.Link)
			if err != nil {
				log.Println(err)
				return *nd, err
			}
			nd.Name = tuic.Name
		}
	}
	return *nd, nil
}
func NodeUpdadte(c *gin.Context) {
	// var node models.Node
	NewName := c.PostForm("name")
	Newlink := c.PostForm("link")
	id := c.PostForm("id")
	group := c.PostForm("group")        // 分组
	groups := strings.Split(group, ",") // 分组列表
	index, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "id 不能为空或者格式不正确",
		})
		return

	}
	if NewName == "" || Newlink == "" {
		c.JSON(400, gin.H{
			"msg": "节点名称 or 备注不能为空",
		})
		return
	}
	OldNode := &models.Node{
		ID: index,
	}
	NewNode := &models.Node{
		Name: NewName,
		Link: Newlink,
	}
	var gns []models.GroupNode
	if groups != nil || len(groups) > 0 {
		for _, g := range groups {
			TempGn := models.GroupNode{
				Name: strings.TrimSpace(g), // 去除分组名称两端空格
			}
			gns = append(gns, TempGn) // 生成分组列表
		}

	}
	err = OldNode.UpdateGroup(gns) // 更新分组
	if err != nil {
		c.JSON(400, gin.H{
			"msg": fmt.Sprintf("更新失败: %s", err.Error()),
		})
		return
	}
	err = OldNode.UpdateNode(NewNode)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": fmt.Sprintf("更新失败: %s", err.Error()),
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
	var ns []models.Node
	ns, err := models.GetNodeList()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "node list error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": ns,
		"msg":  "node get",
	})
}

// 获取分组列表
func GroupNodeGet(c *gin.Context) {
	var Gns []models.GroupNode
	Gns, err := models.GetGroupNodeList()
	var data []string
	for _, g := range Gns {
		data = append(data, g.Name)
	}
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": data,
		"msg":  "GroupNode get",
	})
}

// 设置关联分组
func GroupNodeSet(c *gin.Context) {
	// var n models.Node
	var gns []models.GroupNode
	var FirstGroup models.GroupNode
	name := c.PostForm("name")
	group := c.PostForm("group")

	// 将group分割成多个分组
	groups := strings.Split(group, ",")
	if len(groups) == 0 {
		c.JSON(400, gin.H{
			"msg": "分组不能为空",
		})
		return
	}
	log.Println("分组列表:", groups, "数组长度", len(groups))

	// 循环生成或绑定分组
	for _, g := range groups {
		// 如果group为空，跳过
		if strings.TrimSpace(g) == "" {
			log.Println("分组名为空，跳过")
			continue
		}
		log.Println("分组名:", g)
		FirstGroup.Name = g
		err := FirstGroup.Add()
		if err != nil {
			log.Println("添加分组失败:", err)
			c.JSON(400, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// 查找分组并将数据FirstGroup填充 并且插入给gns
		result := models.DB.Model(models.GroupNode{}).Where("name = ?", g).First(&FirstGroup)
		log.Println("FirstGroup", FirstGroup)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println(result.Error)
			c.JSON(400, gin.H{
				"msg": result.Error,
			})
			return
		}
		gns = append(gns, FirstGroup)
	}

	n := models.Node{Name: name}
	err := n.UpdateGroup(gns)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "更新关联分组成功",
	})
}

// 添加节点
func NodeAdd(c *gin.Context) {
	var n models.Node
	link := c.PostForm("link")
	name := c.PostForm("name")
	group := c.PostForm("group")
	n = models.Node{
		Name: name,
		Link: link,
	}
	if link == "" && !strings.Contains(link, "://") {
		c.JSON(400, gin.H{
			"msg": "link不能为空或者格式不正确,请检查链接是否包含协议头,例如 http:// 或 https://",
		})
		return
	}
	// 解码节点名称
	n, err := DocodeNodeName(&n)
	if err != nil {
		log.Println("解码节点名称错误:", err)
		c.JSON(400, gin.H{
			"msg": "解码节点名称错误",
		})
		return
	}

	// 添加节点
	err = n.Add()
	if err != nil {
		log.Println("添加节点失败:", err)
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 关联分组开始
	if strings.TrimSpace(group) != "" { // 去除空格后判断分组是否为空
		groups := strings.Split(group, ",") // 允许多个分组用逗号分隔
		if groups != nil || len(groups) > 0 {
			for _, g := range groups {
				gn := &models.GroupNode{Name: g}
				err = gn.Add()
				if err != nil {
					// 分组不存在
					log.Println(err)
					c.JSON(400, gin.H{
						"msg": err,
					})
					return
				}
				// 分组存在，关联节点
				if err := gn.Ass(&n); err != nil {
					log.Println("关联失败:", err)
					c.JSON(400, gin.H{
						"msg": err,
					})
					return
				}

			}
		}
	}
	//关联分组结束

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "添加成功",
	})
}

// 删除节点
func NodeDel(c *gin.Context) {
	var n models.Node
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, _ := strconv.Atoi(id)
	n.ID = x
	err := n.Del()
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

// 删除分组
func NodesGroup(c *gin.Context) {
	var gn models.GroupNode
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, _ := strconv.Atoi(id)
	gn.ID = x
	err := gn.Del()
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
	var nodes []models.Node
	nodes, err := models.GetNodeList()
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

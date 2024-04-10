package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sublink/models"
	"sublink/node"

	"github.com/gin-gonic/gin"
)

func GetV2ray(c *gin.Context) {
	var sub models.Subcription
	subname := c.Param("subname")
	subname = node.Base64Decode(subname)
	sub.Name = subname
	err := sub.Find()
	if err != nil {
		c.Writer.WriteString("找不到这个订阅:" + subname)
		return
	}
	err = sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取错误")
		return
	}
	baselist := ""
	// for _, v := range sub.Nodes {
	// 	baselist += v.Link + "\n"
	// }
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			baselist += strings.Join(links, "\n") + "\n"
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := node.Base64Decode(string(body))
			baselist += nodes + "\n"
		// 默认
		default:
			baselist += v.Link + "\n"
		}
	}
	Content_Disposition := fmt.Sprintf("inline; filename=%s.txt", subname)
	c.Set("subname", subname)
	c.Writer.Header().Set("Content-Disposition", Content_Disposition)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(node.Base64Encode(baselist))
}
func GetClash(c *gin.Context) {
	var sub models.Subcription
	subname := c.Param("subname")
	subname = node.Base64Decode(subname)
	sub.Name = subname
	err := sub.Find()
	if err != nil {
		c.Writer.WriteString("找不到这个订阅:" + subname)
		return
	}
	err = sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取错误")
		return
	}
	urls := []string{}
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			urls = append(urls, links...)
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := node.Base64Decode(string(body))
			links := strings.Split(nodes, "\n")
			urls = append(urls, links...)
		// 默认
		default:
			urls = append(urls, v.Link)
		}
	}

	var configs node.SqlConfig
	err = json.Unmarshal([]byte(sub.Config), &configs)
	if err != nil {
		c.Writer.WriteString("配置读取错误")
		return
	}
	DecodeClash, err := node.EncodeClash(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", subname)
	c.Writer.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	c.Writer.WriteString(string(DecodeClash))
}

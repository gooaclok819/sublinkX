package api

import (
	"encoding/json"
	"fmt"
	"log"
	"sublink/models"
	"sublink/node"

	"github.com/gin-gonic/gin"
)

func GetV2ray(c *gin.Context) {
	var sub models.Subcription
	subname := c.Param("subname")
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
	for _, v := range sub.Nodes {
		baselist += v.Link + "\n"
	}
	Content_Disposition := fmt.Sprintf("attachment; filename=%s.txt", subname)
	c.Writer.Header().Set("Content-Disposition", Content_Disposition)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(node.Base64Encode(baselist))
}
func GetClash(c *gin.Context) {
	var sub models.Subcription
	subname := c.Param("subname")
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
		urls = append(urls, v.Link)
	}

	var configs node.SqlConfig
	log.Println(sub.Config)
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
	c.Set("name", subname)
	c.Writer.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	c.Writer.WriteString(string(DecodeClash))
}

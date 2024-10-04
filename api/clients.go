package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	c.Set("subname", subname)
	filename := fmt.Sprintf("%s.txt", subname)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
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
	filename := fmt.Sprintf("%s.yaml", subname)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteString(string(DecodeClash))
}
func GetSurge(c *gin.Context) {
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
	// log.Println("surge路径:", configs)
	DecodeClash, err := node.EncodeSurge(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", subname)
	filename := fmt.Sprintf("%s.conf", subname)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	host := c.Request.Host
	url := c.Request.URL.String()
	// 如果包含头部更新信息
	if strings.Contains(DecodeClash, "#!MANAGED-CONFIG") {
		c.Writer.WriteString(DecodeClash)
		return
	}
	// 否则就插入头部更新信息
	interval := fmt.Sprintf("#!MANAGED-CONFIG %s interval=86400 strict=false", host+url)
	c.Writer.WriteString(string(interval + "\n" + DecodeClash))
}

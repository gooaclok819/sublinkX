package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sublink/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func GetIp(c *gin.Context) {
	c.Next()
	func() {
		subname, _ := c.Get("subname")

		ip := c.ClientIP()
		resp, err := http.Get(fmt.Sprintf("https://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true", ip))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		utf8Body, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body)
		type IpInfo struct {
			Addr string `json:"addr"`
			Ip   string `json:"ip"`
		}
		ipinfo := IpInfo{}
		err = json.Unmarshal(utf8Body, &ipinfo)
		if err != nil {
			log.Println(err)
			return
		}
		var sub models.Subcription
		if subname, ok := subname.(string); ok {
			sub.Name = subname
		}
		err = sub.Find()
		if err != nil {
			log.Println(err)
			return
		}
		var iplog models.SubLogs
		iplog.IP = ip
		err = iplog.Find(sub.ID)
		// 如果没有找到记录
		if err != nil {
			iploga := []models.SubLogs{
				{IP: ip,
					Addr:          ipinfo.Addr,
					SubcriptionID: sub.ID,
					Date:          time.Now().Format("2006-01-02 15:04:05"),
					Count:         1,
				},
			}
			sub.SubLogs = iploga
			err = sub.Update()
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			// 更新访问次数
			iplog.Count++
			iplog.Date = time.Now().Format("2006-01-02 15:04:05")
			err = iplog.Update()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

}

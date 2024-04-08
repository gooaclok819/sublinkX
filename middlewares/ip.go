package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sublink/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func GetIp(c *gin.Context) {
	c.Next()
	go func() {
		name, _ := c.Get("name")
		ip := c.ClientIP()
		resp, err := http.Get(fmt.Sprintf("https://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true", ip))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
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
		sub.Name = name.(string)
		err = sub.Find()
		if err != nil {
			log.Println(err)
			return
		}
		var iplog models.SubLogs
		iplog.IP = ip
		err = iplog.Find()
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
			err = iplog.Update()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

}

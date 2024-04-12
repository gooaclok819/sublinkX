package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Temp struct {
	File       string `json:"file"`
	Text       string `json:"text"`
	CreateDate string `json:"create_date"`
}

func GetTempS(c *gin.Context) {
	files, err := os.ReadDir("./template")
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var temps []Temp
	for _, file := range files {
		info, _ := file.Info()
		time := info.ModTime().Format("2006-01-02 15:04:05")
		text, _ := os.ReadFile("./template/" + file.Name())
		temp := Temp{
			File:       file.Name(),
			Text:       string(text),
			CreateDate: time,
		}
		temps = append(temps, temp)
	}
	if len(temps) == 0 {
		c.JSON(200, gin.H{
			"code": "00000",
			"data": []string{},
			"msg":  "ok",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": temps,
		"msg":  "ok",
	})
}
func UpdateTemp(c *gin.Context) {
	filename := c.PostForm("filename")
	oldname := c.PostForm("oldname")
	text := c.PostForm("text")
	err := os.Rename("./template/"+oldname, "./template/"+filename)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "改名失败",
		})
		return
	}
	err = os.WriteFile("./template/"+filename, []byte(text), 0666)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "修改失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "修改成功",
	})

}
func AddTemp(c *gin.Context) {
	filename := c.PostForm("filename")
	text := c.PostForm("text")
	if filename == "" || text == "" {
		c.JSON(400, gin.H{
			"msg": "文件名或者类型或内容不能为空",
		})
		return
	}
	// 检查文件是否存在
	_, err := os.ReadFile("./template/" + filename)
	if err == nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "文件已存在",
		})
		return
	}
	// 检查目录是否创建
	_, err = os.Stat("./template/")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./template/", os.ModePerm)
		}
	}
	err = os.WriteFile("./template/"+filename, []byte(text), 0666)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "上传失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "上传成功",
	})
}
func DelTemp(c *gin.Context) {
	filename := c.PostForm("filename")
	_, err := os.ReadFile("./template/" + filename)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"msg": "文件不存在",
		})
		return
	}
	err = os.Remove("./template/" + filename)
	if err != nil {
		log.Println(err)
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

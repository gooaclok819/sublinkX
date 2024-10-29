package routers

import "github.com/gin-gonic/gin"

func Version(r *gin.Engine, version string) {

	r.GET("/api/v1/version", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"code": "00000",
			"data": version,
		})
	})
}

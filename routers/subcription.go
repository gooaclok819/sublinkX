package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func Subcription(r *gin.Engine) {
	SubcriptionGroup := r.Group("/api/v1/subcription")
	{
		SubcriptionGroup.POST("/add", api.SubAdd)
		SubcriptionGroup.DELETE("/delete", api.SubDel)
		SubcriptionGroup.GET("/get", api.SubGet)
		SubcriptionGroup.POST("/update", api.SubUpdate)
		SubcriptionGroup.POST("/sort", api.SubUpdateSort)
	}
}

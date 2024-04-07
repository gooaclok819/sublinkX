package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func Mentus(r *gin.Engine) {
	MentusGroup := r.Group("/api/v1/menus")
	{
		// MentusGroup.GET("/menus", api.GetMenus)
		MentusGroup.GET("/routes", api.GetMenus)
	}
}

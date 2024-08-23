package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func Subscription(r *gin.Engine) {
	SubscriptionGroup := r.Group("/api/v1/subscription")
	{
		SubscriptionGroup.POST("/add", api.SubAdd)
		SubscriptionGroup.DELETE("/delete", api.SubDel)
		SubscriptionGroup.GET("/get", api.SubGet)
		SubscriptionGroup.POST("/update", api.SubUpdate)
	}

}

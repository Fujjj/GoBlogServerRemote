package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type FriendLinkRouter struct{}

func (f *FriendLinkRouter) InitFriendLinkRouter(AdminRouter *gin.RouterGroup, PublicGroup *gin.RouterGroup) {
	fAdminRouter := AdminRouter.Group("friendLink")
	fPublicRouter := PublicGroup.Group("friendLink")

	friendLinkApi := api.ApiGroupApp.FriendLinkApi

	{
		fAdminRouter.POST("create", friendLinkApi.FriendLinkCreate)
		fAdminRouter.DELETE("delete", friendLinkApi.FriendLinkDelete)
		fAdminRouter.PUT("update", friendLinkApi.FriendLinkUpdate)
		fAdminRouter.GET("list", friendLinkApi.FriendLinkList)
	}
	{
		fPublicRouter.GET("info", friendLinkApi.FriendLinkInfo)
	}
}

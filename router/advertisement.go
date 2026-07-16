package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type AdvertisementRouter struct {
}

func (a *AdvertisementRouter) InitAdvertisementRouter(AdminRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adAdminRouter := AdminRouter.Group("advertisement")
	adPublicRouter := PublicRouter.Group("advertisement")

	advertisementApi := api.ApiGroupApp.AdvertisementApi
	{
		adAdminRouter.POST("create", advertisementApi.AdvertisementCreate)
		adAdminRouter.DELETE("delete", advertisementApi.AdvertisementDelete)
		adAdminRouter.PUT("update", advertisementApi.AdvertisementUpdate)
		adAdminRouter.GET("list", advertisementApi.AdvertisementList)
	}
	{
		adPublicRouter.GET("info", advertisementApi.AdvertisementInfo)
	}
}

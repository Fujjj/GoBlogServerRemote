package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ConfigRouter struct {
}

func (c *ConfigRouter) InitConfigRouter(Router *gin.RouterGroup) {
	configRouter := Router.Group("config")

	configApi := api.ApiGroupApp.ConfigApi
	{
		//同一个路由不同的方法是合法的
		configRouter.GET("website", configApi.GetWebsite)
		configRouter.PUT("website", configApi.UpdateWebsite)
		configRouter.GET("system", configApi.GetSystem)
		configRouter.PUT("system", configApi.UpdateSystem)
		configRouter.GET("email", configApi.GetEmail)
		configRouter.PUT("email", configApi.UpdateEmail)
		configRouter.GET("qq", configApi.GetQQ)
		configRouter.PUT("qq", configApi.UpdateQQ)
		configRouter.GET("qiniu", configApi.GetQiniu)
		configRouter.PUT("qiniu", configApi.UpdateQiniu)
		configRouter.GET("jwt", configApi.GetJwt)
		configRouter.PUT("jwt", configApi.UpdateJwt)
		configRouter.GET("gaode", configApi.GetGaode)
		configRouter.PUT("gaode", configApi.UpdateGaode)
	}
}

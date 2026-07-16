package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type WebsiteRouter struct {
}

func (websiteRouter *WebsiteRouter) InitWebsiteRouter(AdminRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	websitePublicRouter := PublicRouter.Group("website")
	websiteAdminRouter := AdminRouter.Group("website")

	websiteApi := api.ApiGroupApp.WebsiteApi
	{
		websiteAdminRouter.POST("addCarousel", websiteApi.WebsiteAddCarousel)
		websiteAdminRouter.PUT("cancelCarousel", websiteApi.WebsiteCancelCarousel)
		websiteAdminRouter.POST("createFooterLink", websiteApi.WebsiteCreateFooterLink)
		websiteAdminRouter.DELETE("deleteFooterLink", websiteApi.WebsiteDeleteFooterLink)
	}
	{
		websitePublicRouter.GET("logo", websiteApi.WebsiteLogo)
		websitePublicRouter.GET("title", websiteApi.WebsiteTitle)
		websitePublicRouter.GET("info", websiteApi.WebsiteInfo)
		websitePublicRouter.GET("carousel", websiteApi.WebsiteCarousel)
		// websitePublicRouter.GET("news", websiteApi.WebsiteNews)
		websitePublicRouter.GET("calendar", websiteApi.WebsiteCalendar)
		websitePublicRouter.GET("footerLink", websiteApi.WebsiteFooterLink)
	}

}

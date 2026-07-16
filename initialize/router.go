package initialize

import (
	"net/http"
	"server/global"
	"server/middleware"
	"server/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	//根据配置设置gin模式
	gin.SetMode(global.Configs.System.Env)
	Router := gin.Default()
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	//初始化并启动基于cookie的session中间件
	var store = cookie.NewStore([]byte(global.Configs.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
	// 将指定目录下的文件提供给客户端
	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	//配置静态文件服务，使得上传的文件可以通过 HTTP 直接访问。
	Router.StaticFS(global.Configs.Upload.Path, http.Dir(global.Configs.Upload.Path))

	//创建路由组
	routerGroup := router.RouterGroupApp
	//创建带有统一前缀的路由组
	publicGroup := Router.Group(global.Configs.System.RouterPrefix)
	privateGroup := Router.Group(global.Configs.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuth())
	adminGroup := Router.Group(global.Configs.System.RouterPrefix)
	adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
	{
		//把publicGroup注册在基础路由下
		routerGroup.InitBaseRouter(publicGroup)
	}
	{
		//把privateGroup、publicGroup、adminGroup注册在用户路由下
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
		//把privateGroup、publicGroup、adminGroup注册在文章路由下
		routerGroup.InitArticleRouter(privateGroup, publicGroup, adminGroup)
		//把privateGroup、publicGroup、adminGroup注册在评论路由下
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
		//把privateGroup、publicGroup、adminGroup注册在反馈路由下
		routerGroup.InitFeedbackRouter(privateGroup, publicGroup, adminGroup)
	}
	{
		//把adminGroup注册在图片路由下
		routerGroup.InitImageRouter(adminGroup)
		//将adminGroup、publicGroup注册在广告路由和友链路由下
		routerGroup.InitAdvertisementRouter(adminGroup, publicGroup)
		routerGroup.InitFriendLinkRouter(adminGroup, publicGroup)
		//把adminGroup、publicGroup注册在网站路由下
		routerGroup.InitWebsiteRouter(adminGroup, publicGroup)
		//把adminGroup注册在配置路由下，配置路由用于获取或修改config.yaml中的配置
		routerGroup.InitConfigRouter(adminGroup)
	}

	return Router
}

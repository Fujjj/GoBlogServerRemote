package router

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

// 初始化用户路由组,使用大写的变量名是为了避免与包名重合
func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	userRouter := Router.Group("user")
	userPublicRouter := PublicRouter.Group("user")

	//日志记录的路由
	userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())

	userAdminRouter := AdminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi
	{
		userRouter.POST("logout", userApi.Logout)
		userRouter.PUT("resetPassword", userApi.UserResetPassword)
		userRouter.GET("info", userApi.UserInfo)             // 获取当前登录的用户信息
		userRouter.PUT("changeInfo", userApi.UserChangeInfo) //修改用户信息
		userRouter.GET("weather", userApi.UserWeather)       //获取用户天气
		userRouter.GET("chart", userApi.UserChart)           //获取用户图表数据，包括登录人数和注册人数
	}
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword) //找回密码
		userPublicRouter.GET("card", userApi.UserCard)                  //查看用户卡片信息
	}
	{
		userLoginRouter.POST("register", userApi.Register)
		userLoginRouter.POST("login", userApi.Login)
	}
	{
		userAdminRouter.GET("list", userApi.UserList)           //用户列表
		userAdminRouter.PUT("freeze", userApi.UserFreeze)       //冻结用户
		userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)   //解冻用户
		userAdminRouter.GET("loginList", userApi.UserLoginList) //登录日志列表
	}
}

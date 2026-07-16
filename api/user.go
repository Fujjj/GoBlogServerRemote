package api

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type UserApi struct {
}

func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Successful logout", c)
}

func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//user id 通过后端获取
	req.UserID = utils.GetUserID(c)
	if err := userService.UserResetPassword(req); err != nil {
		global.Log.Error("Failed to reset password:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Password updated successfully,please log in again", c)
	userApi.Logout(c)
}

func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user info:", zap.Error(err))
		response.FailWithMessage("Failed to get user info:", c)
		return
	}
	response.OkWithData(user, c)

}

func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//通过前端传来的请求头中的jwt获取用户id
	req.UserID = utils.GetUserID(c)
	if err := userService.UserChangeInfo(req); err != nil {
		global.Log.Error("Failed to change user info:", zap.Error(err))
		response.FailWithMessage("Failed to change user info:", c)
		return
	}
	response.OkWithMessage("User info changed successfully", c)
}

func (userApi *UserApi) UserWeather(c *gin.Context) {
	weather, err := userService.UserWeather(c)
	if err != nil {
		global.Log.Error("Failed to get weather:", zap.Error(err))
		response.FailWithMessage("Failed to get weather:", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据，登录和注册人数
func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to get user chart:", c)
		return
	}
	response.OkWithData(data, c)
}

func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取会话
	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	// 使用双返回值断言，避免 panic
	emailStr, ok := savedEmail.(string)
	if !ok || emailStr != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	// 使用双返回值断言，避免 panic
	codeStr, ok := savedCode.(string)
	if !ok || codeStr != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	if err := userService.ForgotPassword(req); err != nil {
		global.Log.Error("Failed to update password:", zap.Error(err))
		response.FailWithMessage("Failed to update password", c)
		return
	}
	response.OkWithMessage("Successfully updated password", c)

}

func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get user card:", zap.Error(err))
		response.FailWithMessage("Failed to get user card", c)
		return
	}
	response.OkWithData(userCard, c)
}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取会话
	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	// 使用双返回值断言，避免 panic
	emailStr, ok := savedEmail.(string)
	if !ok || emailStr != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	// 使用双返回值断言，避免 panic
	codeStr, ok := savedCode.(string)
	if !ok || codeStr != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	// 注册成功后，生成 token 并返回
	//实现注册完自动登录
	userApi.TokenNext(c, user)
}

func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	// case "qq":
	// 	userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}
		// 登录成功后生成 token
		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// func (userApi *UserApi) QQLogin(c *gin.Context) {
// }

func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	if user.Freeze {
		response.FailWithMessage("The user has been frozen, contact the administrator", c)
		return
	}
	//创建基础声明
	baseClaims := request.BaseClaims{
		UUID:   user.UUID,
		UserID: user.ID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()
	accessClaims := j.CreatAccessClaim(baseClaims)
	accessToken, err := j.CreatAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get access token:", zap.Error(err))
		response.FailWithMessage("Failed to get access token", c)
		return
	}
	refreshClaims := j.CreatRefreshClaim(baseClaims)
	refreshToken, err := j.CreatRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to get refresh token:", zap.Error(err))
		response.FailWithMessage("Failed to get refresh token", c)
		return
	}
	//是否开启了多地点登录拦截
	if !global.Configs.System.UseMultipoint {
		//将refreshToken保存在cookie中
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))

		//在middleware/login_record.go中根据user_id设置登录记录，所以这里要设置user_id
		c.Set("user_id", user.ID)

		//返回成功响应
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Login successful", c)
		return
	}
	//如果开启了多地点登录拦截，新登录会踢掉旧登录。
	// 从redis中获取用户旧的jwt
	if jwtStr, err := jwtService.GetRedisJWT(c, user.UUID); errors.Is(err, redis.Nil) {
		//若redis中无记录，则创建新的jwt，并写入redis中
		if err := jwtService.SetRedisJWT(c, accessToken, user.UUID); err != nil {
			global.Log.Error("Failed to set redis jwt:", zap.Error(err))
			response.FailWithMessage("Failed to set redis jwt", c)
			return
		}
		//将refreshToken保存在cookie中
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))

		//在middleware/login_record.go中根据user_id设置登录记录，所以这里要设置user_id
		c.Set("user_id", user.ID)

		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Login successful", c)
	} else if err != nil { //其他错误，进行错误处理
		global.Log.Error("Failed to get redis jwt:", zap.Error(err))
		response.FailWithMessage("Failed to get redis jwt", c)
		return
	} else {
		// Redis 中已存在该用户的 JWT，将旧的 JWT 加入黑名单，并设置新的 token
		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := jwtService.JoinInBlacklist(blacklist); err != nil {
			global.Log.Error("Failed to join in blacklist:", zap.Error(err))
			response.FailWithMessage("Failed to join in blacklist", c)
			return
		}
		//如果成功加入黑名单，设置新的accessToken进redis
		if err := jwtService.SetRedisJWT(c, accessToken, user.UUID); err != nil {
			global.Log.Error("Failed to set redis jwt:", zap.Error(err))
			response.FailWithMessage("Failed to set redis jwt", c)
			return
		}
		// 设置新的 refreshToken 到 Redis
		if err := jwtService.SetRedisJWT(c, refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	}

}

// UserList 获取用户列表
func (userApi *UserApi) UserList(c *gin.Context) {
	var pageinfo request.UserList
	if err := c.ShouldBindQuery(&pageinfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.UserList(pageinfo)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)

}

func (userApi *UserApi) UserFreeze(c *gin.Context) {
	//通过id冻结用户,管理员可以直接拿到id
	var req request.UserOperation

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := userService.UserFreeze(c, req); err != nil {
		global.Log.Error("Failed to freeze user:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("User frozen successfully", c)

}

func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := userService.UserUnfreeze(req); err != nil {
		global.Log.Error("Failed to unfreeze user:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("User unfrozen successfully", c)
}

// // UserLoginList 获取登录日志列表
func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)

}

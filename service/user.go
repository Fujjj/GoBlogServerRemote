package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

func (userService *UserService) Logout(c *gin.Context) {
	//因为在middleware/jwt.go中JWTAuth()中间件在验证通过后会在将Claims信息存入Context，供后续使用，所以此处能够访问
	uuid := utils.GetUUID(c)
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	global.Redis.Del(c, uuid.String())
	_ = ServiceGroupApp.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})

}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	// 验证旧密码
	if !utils.BcryptCheck(req.Password, user.Password) {
		return errors.New("original password does not match the current account")
	}
	//如果旧密码验证成功，加密新密码
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	err := global.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	//根据req结构体中的信息更新user表
	return global.DB.Model(&user).Updates(req).Error
}

func (userService *UserService) UserWeather(c *gin.Context) (string, error) {
	ip := c.ClientIP()

	//防止每一次都调用天气接口浪费性能和额度，先从redis中获取天气信息
	result, err := global.Redis.Get(c, "weather-"+ip).Result()
	if err != nil { //redis中不存在
		//获取城市编码
		ipResponse, err := ServiceGroupApp.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}

		//拼接字符串
		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"
		//将天气数据存储到redis中
		if err := global.Redis.Set(c, "weather-"+ip, weather, 3*time.Hour).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	//如果redis中有数据
	return result, nil

}

func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件，用于筛选出最近 N 天内创建的数据记录。
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var res response.UserChart

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}

	return res, nil
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid=?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil

}

func (userService *UserService) Register(u database.User) (database.User, error) {
	//如果该邮箱已存在，则返回错误
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}

	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		//校验密码是否正确
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("password error")
		}
		return user, nil
	}
	return database.User{}, errors.New("email address does not exist")
}

func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB
	//若uuid存在则查询该用户
	if info.UUID != nil {
		db = db.Where("uuid = ?", *info.UUID)
	}

	//若注册方式存在则查询该注册方式
	if info.Register != nil {
		db = db.Where("register = ?", *info.Register)
	}
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)
}

func (userService *UserService) UserFreeze(c *gin.Context, req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}

	//使当前 JWT Token 失效（踢出登录）
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(c, user.UUID)
	if jwtStr != "" {
		err := ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
		if err != nil {
			return err
		}
	}
	return nil
}

func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}

// 登录日志通过UUID查询,但是database/login的表中存储的是userid
func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	//若存在UUID则转化为userid
	if info.UUID != nil {
		var userID uint
		//根据uuid查询userid，因为uuid是唯一的，只会查到该uuid对应的userid
		if err := global.DB.Model(database.User{}).Where("uuid = ?", *info.UUID).Pluck("id", &userID); err != nil {
			return nil, 0, nil
		}
		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option)
}

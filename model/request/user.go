package request

type Register struct {
	Username         string `json:"username" binding:"required,max=20"`
	Password         string `json:"password" binding:"required,min=8,max=16"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
}

type Login struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=16"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

type ForgotPassword struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
	NewPassword      string `json:"new_password" binding:"required,min=8,max=16"`
}

type UserCard struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
}

type UserResetPassword struct {
	UserID      uint   `json:"-"` //不应该从前端传进来，而是调用后端的方法从请求头中获取
	Password    string `json:"password" binding:"required,min=8,max=16"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=16"`
}

type UserChangeInfo struct {
	UserID    uint   `json:"-"`
	Username  string `json:"username" binding:"required,max=20"`
	Address   string `json:"address" binding:"max=200"`
	Signature string `json:"signature" binding:"max=320"`
}

type UserChart struct {
	Date int `json:"date" form:"date" binding:"required,oneof=7 30 90 180 365"`
}

// 用户列表 根据uuid进行查询，Register是注册来源
type UserList struct {
	UUID     *string `json:"uuid" form:"uuid"` //*代表可传可不传
	Register *string `json:"register" form:"register"`
	PageInfo
}

type UserOperation struct {
	ID uint `json:"id" binding:"required"`
}

type UserLoginList struct {
	UUID *string `json:"uuid" form:"uuid"`
	PageInfo
}

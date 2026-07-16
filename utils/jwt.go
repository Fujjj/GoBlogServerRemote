package utils

import (
	"errors" // jwt-go v3
	"server/global"
	"server/model/request"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
}

// 自定义错误
var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

// 返回一个JWT实例
func NewJWT() *JWT {
	return &JWT{
		AccessTokenSecret:  []byte(global.Configs.Jwt.AccessTokenSecret),
		RefreshTokenSecret: []byte(global.Configs.Jwt.RefreshTokenSecret),
	}
}

// 构建访问声明
func (j *JWT) CreatAccessClaim(baseClaims request.BaseClaims) request.JwtCustomClaims {
	ept, _ := ParseDuration(global.Configs.Jwt.AccessTokenExpiryTime)
	claims := request.JwtCustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go Blog"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ept)),
			Issuer:    global.Configs.Jwt.Issuer,
		},
	}
	return claims
}

// 构建访问令牌
func (j *JWT) CreatAccessToken(claims request.JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessTokenSecret)
}

// 创建刷新声明
func (j *JWT) CreatRefreshClaim(baseClaims request.BaseClaims) request.JwtCustomRefreshClaims {
	ept, _ := ParseDuration(global.Configs.Jwt.AccessTokenExpiryTime)
	claims := request.JwtCustomRefreshClaims{
		UserID: baseClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go_Blog"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ept)),
			Issuer:    global.Configs.Jwt.Issuer,
		},
	}
	return claims
}

// 创建刷新令牌
func (j *JWT) CreatRefreshToken(claims request.JwtCustomRefreshClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

// ParseAccessToken 解析 Access Token，验证 Token 并返回 Claims 信息
// 返回指针类型的Claims是因为ParseToken函数会将解析出的JSON数据unmarshal到Claims中
func (j *JWT) ParseAccessToken(tokenString string) (*request.JwtCustomClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomClaims{}, j.AccessTokenSecret) // 解析 Token
	if err != nil {
		return nil, err
	}
	if customClaims, ok := claims.(*request.JwtCustomClaims); ok { // 确保解析出的 Claims 类型正确
		return customClaims, nil
	}
	return nil, ErrTokenInvalid // 如果解析结果无效，返回 TokenInvalid 错误
}

// ParseRefreshToken 解析 Refresh Token，验证 Token 并返回 Claims 信息
func (j *JWT) ParseRefreshToken(tokenString string) (*request.JwtCustomRefreshClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomRefreshClaims{}, j.RefreshTokenSecret) // 解析 Token
	if err != nil {
		return nil, err
	}
	if refreshClaims, ok := claims.(*request.JwtCustomRefreshClaims); ok { // 确保解析出的 Claims 类型正确
		return refreshClaims, nil
	}
	return nil, ErrTokenInvalid // 如果解析结果无效，返回 TokenInvalid 错误
}

// parseToken 通用的 Token 解析方法，验证 Token 是否有效并返回 Claims
// parseToken返回的interface{}内是解析好的claims。而回调函数的参数interface{}内是密钥本身（secretKey），它返回的 error 是指“获取密钥过程中是否出错”（比如找不到密钥），而不是“Token 是否有效”
// token是锁，回调函数是钥匙管理员，锁住的东西是claims
func (j *JWT) parseToken(tokenString string, claims jwt.Claims, secretKey interface{}) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // 返回密钥以验证 Token
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok { // 处理 Token 验证错误
			switch {
			//使用位运算（Bitwise Operation）来检查 JWT 验证错误的具体类型。
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, ErrTokenMalformed // Token 格式错误
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired // Token 已过期
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotValidYet // Token 还未生效
			default:
				return nil, ErrTokenInvalid // 其他错误返回 Token 无效
			}
		}
		return nil, ErrTokenInvalid // 默认返回 Token 无效错误
	}

	if token.Valid { // 如果 Token 验证通过，返回 Claims
		return token.Claims, nil
	}
	return nil, ErrTokenInvalid // Token 无效，返回错误
}

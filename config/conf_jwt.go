package config

type Jwt struct {
	AccessTokenSecret      string `json:"access_token_secret" yaml:"access_token_secret"`             // the key used for generating and verifying access tokens
	RefreshTokenSecret     string `json:"refresh_token_secret" yaml:"refresh_token_secret"`           // used for generating and verifying refresh tokens
	AccessTokenExpiryTime  string `json:"access_token_expiry_time" yaml:"access_token_expiry_time"`   // 访问令牌的过期时间，例如 "15m" 表示 15 分钟
	RefreshTokenExpiryTime string `json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"` // 刷新令牌的过期时间，例如 "30d" 表示 30 天
	Issuer                 string `json:"issuer" yaml:"issuer"`                                       // JWT 的签发者信息，通常是应用或服务的名称
}

package config

// Gaode config，details: https://lbs.amap.com/
type Gaode struct {
	Enable bool   `json:"enable" yaml:"enable"` // whether to use gaode service
	Key    string `json:"key" yaml:"key"`       // the application key of gaode service,used for anthentication and accessing
}

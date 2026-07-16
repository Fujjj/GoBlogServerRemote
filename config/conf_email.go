package config

type Email struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	From     string `json:"from" yaml:"from"`
	NikeName string `json:"nike_name" yaml:"nike_name"`
	Secret   string `json:"secret" yaml:"secret"`
	IsSSL    bool   `json:"is_ssl" yaml:"is_ssl"` //whether ssl encryption is used
}

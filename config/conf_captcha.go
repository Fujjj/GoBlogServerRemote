package config

type Captcha struct {
	Height   int     `json:"height" yaml:"height"`
	Width    int     `json:"width" yaml:"width"`
	Length   int     `json:"length" yaml:"length"`
	MaxSkew  float64 `json:"max_skew" yaml:"max_skew"`   //the maxmium skew for each character
	DotCount int     `json:"dot_count" yaml:"dot_count"` //the number of background dots
}

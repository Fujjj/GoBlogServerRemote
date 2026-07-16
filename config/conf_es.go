package config

// ES ElasticSearch config
type ES struct {
	URL            string `json:"url" yaml:"url"` //the url of Elasticsearch,such as http://localhost:9200
	Username       string `json:"username" yaml:"username"`
	Password       string `json:"password" yaml:"password"`
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` // whether to print Elasticsearch on the console
}

package core

//负责对config.yaml实现序列化和反序列化，需要引入相关的包go get "gopkg.in/yaml.v3"
import (
	"log"
	"server/config"
	"server/utils"

	"gopkg.in/yaml.v3"
)

// InitConf 从 YAML 文件加载配置
func InitConf() *config.Config {
	//获取一个Config实例
	c := &config.Config{}
	//写一个工具加载yaml文件
	yamlConf, err := utils.LoadYaml()
	if err != nil {
		log.Fatalf("Failed to load configuration:%v", err)
	}
	//yamlConf读取的原始字节数据（[]byte），仅为文本内容的二进制表示。yaml.Unmarshal 负责解析字节流，将 YAML 文件中的键值对映射到 c 的对应字段中。
	if err := yaml.Unmarshal(yamlConf, c); err != nil {
		log.Fatalf("Failed to unmarshal yaml:%v", err)
	}
	return c

}

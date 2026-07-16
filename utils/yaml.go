package utils

import (
	"io/fs"
	"os"
	"server/global"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

// LoadYAML 从文件中读取 YAML 数据并返回字节数组
func LoadYaml() ([]byte, error) {
	return os.ReadFile(configFile)
}

// SaveYAML 将全局配置对象保存为 YAML 格式到文件
func SaveYAML() error {
	byteData, err := yaml.Marshal(global.Configs)
	if err != nil {
		return err
	}
	//fs.ModePerm即0777，所有者、组、其他 都具有读写和执行权限
	return os.WriteFile(configFile, byteData, fs.ModePerm)
}

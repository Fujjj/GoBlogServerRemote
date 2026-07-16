package appTypes

import "encoding/json"

// Storage 图片存储类型
type Storage int

const (
	Local Storage = iota // 本地
	Qiniu                // 七牛云
)

/*
为什么storage 实现了 json.Marshaler 和 json.Unmarshaler 接口
	Storage 底层类型为 int，默认 JSON 编码会输出数字。通过实现 MarshalJSON 方法，
	将其转换为 String 方法返回的语义化字符串（如 "本地"）。
	这样，当 Storage 值被编码为 JSON 时，将返回语义化的字符串，而不是数字。

	实现 UnmarshalJSON 方法后，允许接收字符串形式的输入，并通过 ToStorage 函数将其解析回内部的 Storage 类型。
*/
// MarshalJSON 实现了 json.Marshaler 接口
func (s Storage) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口
func (s *Storage) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = ToStorage(str)
	return nil
}

// String 方法返回 Storage 的字符串表示
func (s Storage) String() string {
	var str string
	switch s {
	case Local:
		str = "本地"
	case Qiniu:
		str = "七牛云"
	default:
		str = "未知存储"
	}
	return str
}

// ToStorage 函数将字符串转换为 Storage
func ToStorage(str string) Storage {
	switch str {
	case "本地":
		return Local
	case "七牛云":
		return Qiniu
	default:
		return -1
	}
}

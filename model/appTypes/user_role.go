package appTypes

// RoleID 用户角色
type RoleID int

const (
	Guest RoleID = iota //游客
	User                // 普通用户
	Admin               // 管理员
)

//如果希望 RoleID 在 JSON 中以字符串形式出现（如 "admin"），
// 则需要实现 json.Marshaler 和 json.Unmarshaler 接口：

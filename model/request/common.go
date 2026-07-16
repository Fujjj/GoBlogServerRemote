package request

type PageInfo struct {
	Page     int `json:"page" form:"page"`           //例如第一页
	PageSize int `json:"page_size" form:"page_size"` //例如10条每页
}

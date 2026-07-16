package request

//通过id删除图片
type ImageDelete struct {
	IDs []uint `json:"ids"`
}

type ImageList struct {
	Name     *string `json:"name" form:"name"`
	Category *string `json:"category" form:"category"`
	Storage  *string `json:"storage" form:"storage"`
	PageInfo
}

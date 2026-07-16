package response

type PageResult struct {
	List  interface{} `json:"list"`  //具体记录，例如第 1 页的 10 条用户数据
	Total int64       `json:"total"` //总记录数，例如共5条记录
}

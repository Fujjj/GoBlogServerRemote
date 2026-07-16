package hotSearch

// //Source接口，统一不同平台（百度、快手等）的获取热搜行为。
// type Source interface {
// 	GetHotSearchData(maxNum int) (HotSearchData other.HotSearchData, err error)
// }

// func NewSource(sourceStr string) Source {
// 	switch sourceStr {
// 	case "baidu":
// 		return &Baidu{}
// 	case "kuaishou":
// 		return &Kuaishou{}
// 	case "toutiao":
// 		return &Toutiao{}
// 	case "zhihu":
// 		return &Zhihu{}
// 	default:
// 		return nil
// 	}
// }

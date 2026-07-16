package service

// 将具体的业务服务（如 EsService）聚合到 ServiceGroup 中，并通过全局变量 ServiceGroupApp 提供统一的访问入口。
type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	ImageService
	ArticleService
	CommentService
	AdvertisementService
	FriendLinkService
	FeedbackService
	WebsiteService
	HotSearchService
	CalendarService
	ConfigService
}

var ServiceGroupApp = new(ServiceGroup)

package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FriendLinkApi struct {
}

// FriendLinkCreate 创建友链
func (friendLinkApi *FriendLinkApi) FriendLinkCreate(c *gin.Context) {
	var req request.FriendLinkCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := friendLinkService.FriendLinkCreate(req); err != nil {
		global.Log.Error("Failed to create friend link:", zap.Error(err))
		response.FailWithMessage("Failed to create friend link:", c)
		return
	}
	response.OkWithMessage("Successfully created friend link", c)

}

// FriendLinkDelete 删除友链
func (friendLinkApi *FriendLinkApi) FriendLinkDelete(c *gin.Context) {
	var req request.FriendLinkDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := friendLinkService.FriendLinkDelete(req); err != nil {
		global.Log.Error("Failed to Delete friend link:", zap.Error(err))
		response.FailWithMessage("Failed to Delete friend link:", c)
		return
	}
	response.OkWithMessage("Successfully Delete friend link", c)

}

// FriendLinkUpdate 更新友链
func (friendLinkApi *FriendLinkApi) FriendLinkUpdate(c *gin.Context) {
	var req request.FriendLinkUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := friendLinkService.FriendLinkUpdate(req); err != nil {
		global.Log.Error("Failed to Update friend link:", zap.Error(err))
		response.FailWithMessage("Failed to Update friend link:", c)
		return
	}
	response.OkWithMessage("succcesfully Update friend link", c)
}

// FriendLinkList 获取友链列表
func (friendLinkApi *FriendLinkApi) FriendLinkList(c *gin.Context) {
	var pageInfo request.FriendLinkList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := friendLinkService.FriendLinkList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get friend link list:", zap.Error(err))
		response.FailWithMessage("Failed to get friend link list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// FriendLinkInfo 获取友链信息
func (friendLinkApi *FriendLinkApi) FriendLinkInfo(c *gin.Context) {
	list, total, err := friendLinkService.FriendLinkInfo()
	if err != nil {
		global.Log.Error("Failed to get friend link information:", zap.Error(err))
		response.FailWithMessage("Failed to get friend link information", c)
		return
	}
	response.OkWithData(response.FriendLinkInfo{
		List:  list,
		Total: total,
	}, c)
}

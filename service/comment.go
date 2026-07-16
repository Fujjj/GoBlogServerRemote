package service

import (
	"context"
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type CommentService struct{}

func (commentService *CommentService) CommentInfoByArticleId(req request.CommentInfoByArticleId) ([]database.Comment, error) {
	//通过pid去预加载children字段（子评论），再将子评论递归重复先前操作
	var comments []database.Comment
	//查找一级评论
	if err := global.DB.Where("article_id=? and p_id IS NULL", req.ArticleID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avatar,address,signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}

	//遍历一级评论，递归加载子评论
	for i := range comments {
		if err := commentService.LoadChildren(&comments[i]); err != nil {
			return nil, err
		}
	}
	return comments, nil
}
func (commentService *CommentService) CommentCreate(req request.CommentCreate) error {
	//若AfterCreate钩子函数返回错误，则事务回滚，CommentCreate不会被执行
	return global.DB.Create(&database.Comment{
		ArticleID: req.ArticleID,
		Content:   req.Content,
		PID:       req.PID,
		UserUUID:  req.UserUUID,
	}).Error
}
func (commentService *CommentService) CommentDelete(ctx context.Context, currentUser other.CurrentUser, req request.CommentDelete) error {
	//防止所有评论被删除
	if len(req.IDs) == 0 {
		return nil
	}

	return global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//一次性查出所有评论
		var comments []database.Comment
		if err := tx.Where("id IN ?", req.IDs).Find(&comments).Error; err != nil {
			return err
		}

		//权限校验
		for _, comment := range comments {
			if comment.UserUUID != currentUser.UUID && currentUser.RoleID != appTypes.Admin {
				return errors.New("you do not have permission to delete this comment")
			}
		}

		//递归删除
		for _, comment := range comments {
			if err := commentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
				return err
			}
		}
		return nil

	})
}
func (commentService *CommentService) CommentInfo(uuid uuid.UUID) ([]database.Comment, error) {
	var rawComments []database.Comment
	//查找所有的评论
	err := global.DB.Order("id desc").Where("user_uuid=?", uuid).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avatar,address,signature")
	}).Find(&rawComments).Error
	if err != nil {
		return nil, err
	}

	//加载子评论
	for i := range rawComments {
		if err := commentService.LoadChildren(&rawComments[i]); err != nil {
			return nil, err
		}
	}

	//评论去重，如果当前评论的子评论存在为你的评论，就去除子评论
	var comments []database.Comment
	idMap := commentService.FindChildCommentsIDByRootCommentUserUUID(rawComments)
	for i := range rawComments {
		if _, exists := idMap[rawComments[i].ID]; !exists {
			comments = append(comments, rawComments[i])
		}
	}
	return comments, nil
}
func (commentService *CommentService) CommentNew() ([]database.Comment, error) {
	var comments []database.Comment
	if err := global.DB.Order("id desc").Limit(5).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avatar,address,signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
func (commentService *CommentService) CommentList(info request.CommentList) (interface{}, int64, error) {
	db := global.DB

	//根据ArticleID查询
	if info.ArticleID != nil {
		//取指针指向的真实 string 值,否则会传入指针，造成错误
		db = db.Where("article_id = ?", *info.ArticleID)
	}

	//根据uuid查询
	if info.UserUUID != nil {
		db = db.Where("user_uuid = ?", *info.UserUUID)
	}

	//根据内容查询
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.Comment{}, option)
}

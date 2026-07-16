package service

import (
	"server/global"
	"server/model/database"

	"gorm.io/gorm"
)

// LoadChildren加载一个评论下的所有子评论
func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment
	//通过p_id查询Comment表得到p_id的直接子评论，并预加载所需要的用户信息，写入children
	if err := global.DB.Where("p_id=?", comment.ID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avatar,address,signature")
	}).Find(&children).Error; err != nil {
		return err
	}

	//递归查询子评论
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}

	comment.Children = children
	return nil
}

// DeleteCommentAndChildren 递归删除一条评论及其所有子评论（子孙评论）
func (commentService *CommentService) DeleteCommentAndChildren(tx *gorm.DB, commentID uint) error {
	var children []database.Comment
	//找到第一层子评论
	if err := tx.Where("p_id=?", commentID).Find(&children).Error; err != nil {
		return err
	}
	//递归查找子评论的子评论，直到没有子评论
	for _, child := range children {
		if err := commentService.DeleteCommentAndChildren(tx, child.ID); err != nil {
			return err
		}
	}
	//后序遍历，子->自己，先删除子评论再删除根评论，防止外键约束
	if err := tx.Delete(&database.Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}

func (commentService *CommentService) FindChildCommentsIDByRootCommentUserUUID(comments []database.Comment) map[uint]struct{} {
	result := make(map[uint]struct{})

	// 遍历所有根评论
	for _, rootComment := range comments {
		//防止外层循环变量的地址被复用
		root := rootComment

		// 创建一个递归函数来查找与根评论相同 UserUUID 的子评论
		var findChildren func([]database.Comment)

		findChildren = func(children []database.Comment) {
			// 遍历当前子评论
			for _, child := range children {
				// 如果子评论的 UserUUID 与根评论相同，加入结果 map
				if child.UserUUID == root.UserUUID {
					result[child.ID] = struct{}{}
				}
				// 如果有子评论，继续递归
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}
		}

		// 调用递归函数，查找根评论的所有子评论
		findChildren(rootComment.Children)
	}

	return result
}

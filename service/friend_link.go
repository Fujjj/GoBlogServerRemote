package service

import (
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"gorm.io/gorm"
)

type FriendLinkService struct {
}

func (friendLinkService *FriendLinkService) FriendLinkCreate(req request.FriendLinkCreate) error {
	friendLinkToCreate := database.FriendLink{
		Name:        req.Name,
		Link:        req.Link,
		Logo:        req.Logo,
		Description: req.Description,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.ChangeImagesCategory(tx, []string{friendLinkToCreate.Logo}, appTypes.Logo); err != nil {
			return err
		}
		return tx.Create(&friendLinkToCreate).Error
	})

}
func (friendLinkService *FriendLinkService) FriendLinkDelete(req request.FriendLinkDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var friendLinkToDelete database.FriendLink
			if err := tx.Take(&friendLinkToDelete, id).Error; err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, []string{friendLinkToDelete.Logo}); err != nil {
				return err
			}
			if err := tx.Delete(&friendLinkToDelete).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (friendLinkService *FriendLinkService) FriendLinkUpdate(req request.FriendLinkUpdate) error {
	return global.DB.Model(&database.FriendLink{}).Where("id = ?", req.ID).Updates(database.FriendLink{
		Name:        req.Name,
		Link:        req.Link,
		Description: req.Description}).Error
}
func (friendLinkService *FriendLinkService) FriendLinkList(info request.FriendLinkList) (interface{}, int64, error) {
	db := global.DB

	if info.Name != nil {
		db = db.Where("name like ?", "%"+*info.Name+"%")
	}
	if info.Description != nil {
		db = db.Where("description like ?", "%"+*info.Description+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.FriendLink{}, option)
}
func (friendLinkService *FriendLinkService) FriendLinkInfo() (friendLinks []database.FriendLink, total int64, err error) {
	if err = global.DB.Model(&database.FriendLink{}).Count(&total).Find(&friendLinks).Error; err != nil {
		return nil, 0, err
	}
	return friendLinks, total, nil
}

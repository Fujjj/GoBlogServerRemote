package service

import (
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"server/utils/upload"

	"gorm.io/gorm"
)

type ImageService struct {
}

func (imageService *ImageService) ImageUpload(file *multipart.FileHeader) (string, error) {
	oss := upload.NewOss()
	url, filename, err := oss.UploadImage(file)
	if err != nil {
		return "", err
	}

	return url, global.DB.Create(&database.Image{
		Name:     filename,
		URL:      url,
		Category: appTypes.Null,
		Storage:  global.Configs.System.Storage(),
	}).Error
}

func (imageService *ImageService) ImageDelete(req request.ImageDelete) error {
	//若没有指定ids就删除图片会删除所有图片
	if len(req.IDs) == 0 {
		return nil
	}

	var images []database.Image
	//根据IDs获取所有的图片信息并绑定在images中
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}
	for _, image := range images {
		//因为要删除很多图片，所以启动事务
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			//删除本地/七牛文件
			oss := upload.NewOssWithStorage(image.Storage)
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (imageService *ImageService) ImageList(info request.ImageList) (list interface{}, total int64, err error) {
	//构建查询语句
	db := global.DB

	//根据图片名称进行查询
	if info.Name != nil {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}

	//根据图片分类进行查询
	if info.Category != nil {
		category := appTypes.ToCategory(*info.Category)
		db = db.Where("category = ?", category)
	}

	//根据图片存储进行查询
	if info.Storage != nil {
		storage := appTypes.ToStorage(*info.Storage)
		db = db.Where("storage = ?", storage)
	}

	//进行分页查询
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.Image{}, option)
}

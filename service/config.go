package service

import (
	"server/config"
	"server/global"
	"server/model/appTypes"
	"server/utils"

	"gorm.io/gorm"
)

type ConfigService struct {
}

func (configService *ConfigService) UpdateWebsite(website config.Website) error {
	oldArray := []string{
		global.Configs.Website.Logo,
		global.Configs.Website.FullLogo,
		global.Configs.Website.QQImage,
		global.Configs.Website.WechatImage,
	}

	newArray := []string{
		website.Logo,
		website.FullLogo,
		website.QQImage,
		website.WechatImage,
	}

	added, removed := utils.DiffArrays(oldArray, newArray)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.InitImagesCategory(global.DB, removed); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(global.DB, added, appTypes.System); err != nil {
			return err
		}
		global.Configs.Website = website
		if err := utils.SaveYAML(); err != nil {
			return err
		}
		return nil
	})
}

func (configService *ConfigService) UpdateSystem(system config.System) error {
	global.Configs.System.UseMultipoint = system.UseMultipoint
	global.Configs.System.SessionsSecret = system.SessionsSecret
	global.Configs.System.OssType = system.OssType
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateEmail(email config.Email) error {
	global.Configs.Email = email
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQQ(qq config.QQ) error {
	global.Configs.QQ = qq
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQiniu(qiniu config.Qiniu) error {
	global.Configs.Qiniu = qiniu
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateJwt(jwt config.Jwt) error {
	global.Configs.Jwt = jwt
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateGaode(gaode config.Gaode) error {
	global.Configs.Gaode = gaode
	return utils.SaveYAML()
}

package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 实现了OSS的所有接口
type Qiniu struct {
}

// UploadImage 上传图片到七牛云
func (*Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	//检查文件大小是否超过配置限制
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Configs.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Configs.Upload.Size)

	}

	//获取文件扩展名并验证文件类型
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)

	// 检查扩展名是否在白名单 WhiteImageList 中
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	//生成七牛云上传凭证 (Upload Token)
	// PutPolicy 定义了上传的策略，Scope 指定了目标存储空间 (Bucket)
	putPolicy := storage.PutPolicy{Scope: global.Configs.Qiniu.Bucket}
	// 使用 AccessKey 和 SecretKey 创建 MAC 认证对象
	mac := qbox.NewMac(global.Configs.Qiniu.AccessKey, global.Configs.Qiniu.SecretKey)
	// 根据策略生成上传凭证
	upToken := putPolicy.UploadToken(mac)

	//配置七牛云存储区域和传输协议
	cfg := qiniuConfig()
	// 创建表单上传管理器
	formUploader := storage.NewFormUploader(cfg)
	// 用于接收上传返回的结果
	putRet := storage.PutRet{}
	// 上传额外参数，目前为空
	putExtra := storage.PutExtra{Params: map[string]string{}}

	//生成唯一的文件名
	fileKey := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext

	//打开上传的文件流
	data, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer data.Close()

	//调用上传接口
	// &putRet: 接收返回结果
	// upToken: 上传凭证
	// fileKey: 存储在七牛云的文件的文件名
	// data: 文件数据流
	// file.Size: 文件大小
	// &putExtra: 额外参数
	err = formUploader.Put(context.Background(), &putRet, upToken, fileKey, data, file.Size, &putExtra)
	if err != nil {
		return "", "", err
	}

	//返回完整 URL 和文件 Key
	return global.Configs.Qiniu.ImgPath + putRet.Key, putRet.Key, nil
}

// DeleteImage 从七牛云删除图片
func (*Qiniu) DeleteImage(key string) error {
	//创建 MAC 认证对象
	mac := qbox.NewMac(global.Configs.Qiniu.AccessKey, global.Configs.Qiniu.SecretKey)
	//获取存储配置
	cfg := qiniuConfig()
	//创建 bucket 管理器，用于管理存储空间内的文件
	bucketManager := storage.NewBucketManager(mac, cfg)
	//global.Configs.Qiniu.Bucket: 目标存储空间
	// key: 要删除的文件键名
	return bucketManager.Delete(global.Configs.Qiniu.Bucket, key)
}

// qiniuConfig 根据全局配置构建七牛云 SDK 所需的存储配置对象
// 七牛云不同区域的数据中心地址不同，必须正确配置才能上传
func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      global.Configs.Qiniu.UseHTTPS,
		UseCdnDomains: global.Configs.Qiniu.UseCdnDomains,
	}
	switch global.Configs.Qiniu.Zone {
	case "z0", "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "z1", "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "z2", "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "na0", "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	case "ZoneHuadongZheJiang2":
		cfg.Zone = &storage.ZoneHuadongZheJiang2
	}
	return &cfg
}

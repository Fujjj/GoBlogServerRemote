package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

// Local实现了OSS的所有接口
type Local struct {
}

// UploadImage 上传图片到本地服务器
func (*Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	//检查文件大小是否超过配置限制，将字节转换为 MB 进行比较
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Configs.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Configs.Upload.Size)

	}

	//获取文件后缀并验证文件类型
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	// 检查扩展名是否在白名单 WhiteImageList中
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	//生成唯一的文件名，使用 MD5 哈希原始文件名的一部分，加上当前时间戳，防止文件名冲突
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext

	//按照配置文件构建保存路径，并拼接子目录 "/image/"
	path := global.Configs.Upload.Path + "/image/"

	// 确保目录存在，如果不存在则创建
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}
	//拼接完整的文件保存路径
	filepath := path + filename

	//创建目标文件
	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	//打开上传的文件流
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close() // 确保函数退出时关闭文件句柄

	//将上传文件的内容复制到目标文件中
	if _, err = io.Copy(out, f); err != nil {
		return "", "", err
	}

	//返回的路径以 "/" 开头，方便前端直接作为 URL 使用
	return "/" + filepath, filename, nil
}

// DeleteImage 从本地服务器删除图片
func (*Local) DeleteImage(key string) error {
	// 构建完整的文件路径
	path := global.Configs.Upload.Path + "/image/" + key
	return os.Remove(path)
}

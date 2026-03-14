package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile 保存上传的文件
// saveDir: 保存目录（相对于程序运行目录）
// file: 上传的文件
// 返回：文件相对路径，错误
func SaveUploadedFile(saveDir string, file *multipart.FileHeader) (string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败：%v", err)
	}
	defer src.Close()

	// 创建保存目录
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败：%v", err)
	}

	// 生成唯一文件名
	fileExt := filepath.Ext(file.Filename)
	fileName := generateFileName(fileExt)
	filePath := filepath.Join(saveDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败：%v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败：%v", err)
	}

	// 返回相对路径
	return filepath.Join(saveDir, fileName), nil
}

// generateFileName 生成唯一的文件名
func generateFileName(ext string) string {
	timestamp := time.Now().Format("20060102150405")
	random := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%d%d", timestamp, random)))
	return fmt.Sprintf("%s%s", hex.EncodeToString(hash[:]), ext)
}

package utils

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// CheckExist 检测文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// CheckPermission 检测是否具有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return !os.IsPermission(err)
}

// IsNotExistMkDir 不存在则新建目录
func IsNotExistMkDir(src string) error {
	if exist := CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir 新建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// CheckImageSize 检测图片大小
func CheckImageSize(f multipart.File, maxSize int) bool {
	size, err := GetSize(f)
	if err != nil {
		return false
	}

	return size <= maxSize
}

//  CheckImageExt 检测图片拓展名是否合法
func CheckImageExt(fileName string, allExt ...string) bool {
	ext := path.Ext(fileName)
	for _, allowExt := range allExt {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

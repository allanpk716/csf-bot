package folder

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/utils"
	"os"
	"path/filepath"
)

// AppRootDirPath 获取当前程序的根目录
func AppRootDirPath() string {
	return appRootDirPath
}

// CacheRootDirPath 获取缓存目录
func CacheRootDirPath() string {
	if utils.IsDir(cacheDirPath) == false {
		err := os.MkdirAll(cacheDirPath, os.ModePerm)
		if err != nil {
			logger.Panic("Create Cache Dir Error: " + err.Error())
		}
	}
	return cacheDirPath
}

// CacheCaptchaImgDirPath 获取缓存目录
func CacheCaptchaImgDirPath() string {
	if utils.IsDir(cacheCaptchaImgDirPath) == false {
		err := os.MkdirAll(cacheCaptchaImgDirPath, os.ModePerm)
		if err != nil {
			logger.Panic("Create Cache Dir Error: " + err.Error())
		}
	}
	return cacheCaptchaImgDirPath
}

// SensitiveWordDirPath 获取缓存目录
func SensitiveWordDirPath() string {
	if utils.IsDir(sensitiveWordDirPath) == false {
		err := os.MkdirAll(sensitiveWordDirPath, os.ModePerm)
		if err != nil {
			logger.Panic("Create Dir Error: " + err.Error())
		}
	}
	return sensitiveWordDirPath
}

var (
	cacheDirPath           = filepath.Join(appRootDirPath, "cache")
	cacheCaptchaImgDirPath = filepath.Join(cacheDirPath, "captcha_images")
	sensitiveWordDirPath   = filepath.Join(appRootDirPath, "dict")
)

const (
	appRootDirPath = "."
)

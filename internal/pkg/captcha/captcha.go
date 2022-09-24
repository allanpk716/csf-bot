package captcha

import (
	"encoding/base64"
	"fmt"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/folder"
	"github.com/mojocn/base64Captcha"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	store = base64Captcha.DefaultMemStore
)

// GetCaptcha 写入验证码，并且返回验证码id
func GetCaptcha() (string, string, error) {

	var err error
	// 生成默认数字
	driver := base64Captcha.NewDriverDigit(100, 320, 6, 0.7, 80)
	// 生成base64图片
	c := base64Captcha.NewCaptcha(driver, store)
	// 获取
	code, b64s, err := c.Generate()
	if err != nil {
		return "", "", err
	}
	imageUrl := filepath.Join(folder.CacheCaptchaImgDirPath(), fmt.Sprintf("%s.png", code))
	b64s = b64s[22:]
	b64img, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		return "", "", err
	}
	err = ioutil.WriteFile(imageUrl, b64img, os.ModePerm)
	return code, imageUrl, err
}

// VerifyCaptcha 验证验证码是否正确
func VerifyCaptcha(id, digits string) bool {
	if id == "" || digits == "" {
		return false
	}
	verifyRes := store.Verify(id, digits, false)
	if verifyRes {
		store.Verify(id, digits, true)
		return verifyRes
	} else {
		return false
	}
}

package sensitiveword

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/folder"
	"github.com/importcjj/sensitive"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var Filter *sensitive.Filter

// InitSensitiveWord 加载敏感词库
func InitSensitiveWord() {

	sensitiveWordPath := folder.SensitiveWordDirPath()
	Filter = sensitive.New()
	files, err := ioutil.ReadDir(sensitiveWordPath)
	if err != nil {
		logger.Panic("[InitSensitiveWord] load dict err:", err)
	}
	for _, file := range files {
		// 文件名必须是已解密文件
		if strings.Contains(file.Name(), "dec_") == false {
			continue
		}
		sensitiveFile := filepath.Join(sensitiveWordPath, file.Name())
		err = Filter.LoadWordDict(sensitiveFile)
		if err != nil {
			logger.Panic("[InitSensitiveWord] load sensitive file err:", err, ", file:", sensitiveFile)
		}
	}
}

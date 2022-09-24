package settings

import (
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/folder"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/strcut_json"
	"path/filepath"
)

type Settings struct {
	configFPath                 string // 配置文件的路径
	BotToken                    string // Telegram bot token
	ProxyUrl                    string // http://127.0.0.1:1080 or socks5://127.0.0.1:1081
	CaptchaTimeout              int    // 验证超时时间，单位秒
	UserJoinMessageDelAfter     int    // 用户加入群组后，发送的欢迎消息，多少秒后删除
	NumberOfForbiddenWords      int    // 禁止词数量
	BlockTime                   int64  // 封禁时间，单位秒，-1 是永久
	UserJoinMessage             string // 用户加入群组的欢迎信息
	CaptchaMessage              string // 验证码的欢迎信息
	VerificationCompleteMessage string // 验证完成的欢迎信息
	BlockMessage                string // 被封禁的信息
}

func NewSettings() *Settings {

	nowConfigFPath := filepath.Join(folder.AppRootDirPath(), configName)

	return &Settings{
		configFPath:                 nowConfigFPath,
		BotToken:                    "Telegram bot token",
		ProxyUrl:                    "",
		CaptchaTimeout:              120,
		UserJoinMessageDelAfter:     60,
		NumberOfForbiddenWords:      3,
		UserJoinMessage:             userJoinMessage,
		CaptchaMessage:              captchaMessage,
		VerificationCompleteMessage: verificationCompleteMessage,
		BlockMessage:                blockMessage,
	}
}

func (s *Settings) Read() error {
	return strcut_json.ToStruct(s.configFPath, s)
}

func (s *Settings) Save() error {
	return strcut_json.ToFile(s.configFPath, s)
}

func (s *Settings) ConfigFPath() string {
	return s.configFPath
}

const (
	configName                  = "config.json"
	userJoinMessage             = "欢迎 [%s](%s) 加入 %s\n\n⚠️本群已开启新成员验证功能，未通过验证的用户无法发言 \n\n⏱本条消息 %d 秒后自动删除\n\n👇点击下方按钮自助解除禁言"
	captchaMessage              = "欢迎您加入[%s]！\n\n⚠本群已开启新成员验证功能 \n\n👆为了证明您不是机器人，请发送以上图片验证码内容\n\n🤖机器人将自动验证您发送的验证码内容是否正确\n\n⏱本条验证消息有效期[%d]秒"
	verificationCompleteMessage = "恭喜您成功通过[🤖人机验证]，系统已为您解除禁言限制。\n\n如若还是无法发言，请重启telegram客户端"
	blockMessage                = "\\#封禁预警\n[%s](%s) 请注意,您的消息中含有部分违禁词 \n⚠️您已被系统判断为高风险用户，已被封禁\n系统已向超管发送预警信息，若由超管判定为误杀，会及时将您解除封禁。\n您的违禁词包含：%s"
)

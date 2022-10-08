package telegram_bot

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/captcha"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/command"
	sensitiveword "github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/sensitive_word"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/settings"
	cmap "github.com/orcaman/concurrent-map/v2"
	tb "gopkg.in/telebot.v3"
	"log"
	"net/http"
	"net/url"
	"time"
)

type TelegramBot struct {
	bot                     *tb.Bot
	settings                *settings.Settings
	messageTokenMap         cmap.ConcurrentMap[string]
	oneTimePassCodeMap      cmap.ConcurrentMap[string]
	userCaptchaCodeTable    *captcha.CodeTable
	userCaptchaPendingTable *captcha.PendingTable
}

func NewTelegramBot(settings *settings.Settings) *TelegramBot {

	// TelegramBot 的设置
	botSetting := tb.Settings{
		Token:   settings.BotToken,
		Updates: 100,
		Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, context tb.Context) {
			logger.Error(err)
		},
	}
	// 代理
	if settings.ProxyUrl != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(settings.ProxyUrl)
		}
		httpTransport := &http.Transport{
			Proxy: proxy,
		}
		httpClient := &http.Client{
			Transport: httpTransport,
		}
		botSetting.Client = httpClient
	}
	b, err := tb.NewBot(botSetting)
	if err != nil {
		log.Panicln("NewBot error: ", err)
	}

	return &TelegramBot{
		bot:                     b,
		settings:                settings,
		messageTokenMap:         cmap.New[string](),
		oneTimePassCodeMap:      cmap.New[string](),
		userCaptchaCodeTable:    captcha.NewCodeTable(),
		userCaptchaPendingTable: captcha.NewPendingTable(),
	}
}

func (t *TelegramBot) Start() {

	// 加载敏感词
	sensitiveword.InitSensitiveWord()

	// 机器人是否或者的检测
	t.bot.Handle(command.CMD_Ping, func(c tb.Context) error {
		return c.Send("pong")
	})
	t.bot.Handle(tb.OnUserJoined, t.UserJoinGroup)
	t.bot.Handle(tb.OnText, t.onTextMessage)
	t.bot.Handle(tb.OnUserLeft, func(c tb.Context) error {
		return c.Delete()
	})
	// 主要的验证接口
	t.bot.Handle(command.CMD_Start, t.startCaptcha)
	// 手动验证接口，需要被屏蔽使用一个一次性的验证码
	t.bot.Handle(command.CMD_ManulPassVerify, t.manulPassVerify)
	// 提交一个一次性的验证码，给屏蔽用户用
	t.bot.Handle(command.CMD_SetOneTimePassCode, t.setOneTimePassCode, t.isRootMiddleware)

	logger.Infoln("Telegram Bot Start")

	t.bot.Start()
}

// isManage 判断是否为管理员
func (t *TelegramBot) isManage(chat *tb.Chat, userId int64) bool {
	adminList, err := t.bot.AdminsOf(chat)
	if err != nil {
		return false
	}
	for _, member := range adminList {
		if member.User.ID == userId {
			return true
		}
	}
	return false
}

// isRoot 判断是否为超管
func (t *TelegramBot) isRoot(userid int64) bool {
	for _, id := range t.settings.ManageUsers {
		if userid == id {
			return true
		}
	}
	return false
}

// isManageMiddleware 管理员中间件
func (t *TelegramBot) isManageMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if t.isManage(c.Chat(), c.Sender().ID) {
			return next(c)
		}
		return c.Respond(&tb.CallbackResponse{
			Text:      "您未拥有管理员权限，请勿点击！",
			ShowAlert: true,
		})
	}
}

// isRootMiddleware 超管中间件
func (t *TelegramBot) isRootMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if c.Message().Private() == false || t.isRoot(c.Sender().ID) == false {
			return nil
		}
		return next(c)
	}
}

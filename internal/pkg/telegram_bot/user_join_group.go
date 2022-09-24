package telegram_bot

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/captcha"
	uuid "github.com/satori/go.uuid"
	tb "gopkg.in/telebot.v3"
	"time"
)

// UserJoinGroup 用户加群事件
func (t *TelegramBot) UserJoinGroup(c tb.Context) error {
	var err error
	if err = c.Delete(); err != nil {
		logger.Errorln("[UserJoinGroup] delete join message err:", err)
	}
	// 如果是管理员邀请的，直接通过
	if t.isManage(c.Chat(), c.Sender().ID) {
		return nil
	}
	// 如果是应该被限制的情况下，就不要再发送验证消息了，以免被限制的人利用，重新加群通过验证后解除限制
	userRights, err := t.bot.ChatMemberOf(c.Chat(), c.Message().UserJoined)
	if err != nil {
		return err
	}
	if userRights.Role == tb.Restricted {
		// 已经被限制权限了，不要再发送验证消息了
		return nil
	}
	// ban user，先进行用户权限的限制，如果通过之后再解除限制
	err = t.bot.Restrict(c.Chat(), &tb.ChatMember{
		Rights:          tb.NoRights(),
		User:            c.Message().UserJoined,
		RestrictedUntil: tb.Forever(),
	})
	if err != nil {
		logger.Errorln("[UserJoinGroup] ban user err:", err)
		return err
	}

	return t.makeJoinMessageMenu(c)
}

// makeJoinMessageMenu 制作加群验证消息
func (t *TelegramBot) makeJoinMessageMenu(c tb.Context) error {

	userLink := fmt.Sprintf("tg://user?id=%d", c.Message().UserJoined.ID)
	joinMessage := fmt.Sprintf(t.settings.UserJoinMessage,
		c.Message().UserJoined.LastName+c.Message().UserJoined.FirstName,
		userLink,
		c.Chat().Title,
		t.settings.UserJoinMessageDelAfter)

	joinMessageMenu := &tb.ReplyMarkup{ResizeKeyboard: true}
	chatToken := uuid.NewV4().String()
	// 验证的按钮
	doCaptchaBtn := joinMessageMenu.URL("👉🏻点我开始人机验证🤖", fmt.Sprintf("https://t.me/%s?start=%s", t.bot.Me.Username, chatToken))
	joinMessageMenu.Inline(
		joinMessageMenu.Row(doCaptchaBtn),
		// 添加按钮
		//joinMessageMenu.Row(manageBanBtn, managePassBtn),
	)
	// 添加额外的广告链接
	//LoadAdMenuBtn(joinMessageMenu)
	// 发送验证消息给刚入群的用户
	captchaMessage, err := t.bot.Send(c.Chat(), joinMessage, joinMessageMenu, tb.ModeMarkdownV2)
	if err != nil {
		logger.Error("[UserJoinGroup] send join hint message err:", err)
		return err
	}
	// 设置 token 对于验证消息
	t.messageTokenMap.Set(chatToken, fmt.Sprintf("%d|%d|%s", captchaMessage.ID, c.Chat().ID, c.Chat().Title))
	captchaDataVal := &captcha.Pending{
		PendingMessage: captchaMessage,
		UserId:         c.Message().UserJoined.ID,
		GroupId:        c.Chat().ID,
		JoinAt:         time.Now().Unix(),
	}
	captchaDataKey := fmt.Sprintf("%d|%d", captchaMessage.ID, c.Chat().ID)
	t.userCaptchaPendingTable.Set(captchaDataKey, captchaDataVal)

	// TODO 这里的定时器逻辑需要重构，否则多了就会有大量积累
	time.AfterFunc(time.Duration(t.settings.UserJoinMessageDelAfter)*time.Second, func() {
		if err = t.bot.Delete(captchaMessage); err != nil {
			logger.Error("[UserJoinGroup] delete join hint message err:", err)
		}
	})
	time.AfterFunc(time.Hour, func() {
		t.userCaptchaPendingTable.Del(captchaDataKey)
	})

	return err
}

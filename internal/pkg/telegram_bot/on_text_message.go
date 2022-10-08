package telegram_bot

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/captcha"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

// onTextMessage 文本消息
func (t *TelegramBot) onTextMessage(c tb.Context) error {
	// 私聊走入群验证操作
	if c.Message().Private() {
		return t.verificationProcess(c)
	}
	// 否则走广告阻止监听
	//logger.Infoln(c.Chat().Title, "GroupID:", c.Chat().ID)

	return t.adBlock(c)
}

// verificationProcess 验证处理
func (t *TelegramBot) verificationProcess(c tb.Context) error {

	userIdStr := strconv.FormatInt(c.Sender().ID, 10)
	captchaCode, bok := t.userCaptchaCodeTable.Get(userIdStr)
	if bok == false || captchaCode == nil || captchaCode.UserId != c.Sender().ID {
		return nil
	}
	// 验证
	replyCode := c.Message().Text
	if captcha.VerifyCaptcha(captchaCode.Code, replyCode) == false {
		return nil
	}
	// 解禁用户
	err := t.bot.Restrict(&tb.Chat{ID: captchaCode.GroupId}, &tb.ChatMember{
		User:   &tb.User{ID: captchaCode.UserId},
		Rights: tb.NoRestrictions(),
	})
	if err != nil {
		logger.Error("[onTextMessage] unban err:", err)
		return c.Send("服务器异常~，请稍后重试~")
	}
	t.userCaptchaCodeTable.Del(userIdStr)
	t.userCaptchaPendingTable.Del(fmt.Sprintf("%d|%d", captchaCode.PendingMessage.ID, captchaCode.PendingMessage.Chat.ID))
	//删除验证消息
	if err = t.bot.Delete(captchaCode.CaptchaMessage); err != nil {
		logger.Error("[onTextMessage] delete captcha message err:", err)
	}
	if err = t.bot.Delete(captchaCode.PendingMessage); err != nil {
		logger.Error("[onTextMessage] delete pending message err:", err)
	}
	return c.Send(t.settings.VerificationCompleteMessage)
}

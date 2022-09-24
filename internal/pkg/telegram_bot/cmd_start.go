package telegram_bot

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/captcha"
	tb "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"strings"
	"time"
)

// startCaptcha 开始验证
func (t *TelegramBot) startCaptcha(c tb.Context) error {

	chatToken := c.Message().Payload
	// 不是私聊或者载荷为空
	if !c.Message().Private() || chatToken == "" {
		return nil
	}
	payload, ok := t.messageTokenMap.Get(chatToken)
	if ok == false {
		return nil
	}
	// payload 不能正常解析
	payloadSlice := strings.Split(payload, "|")
	if len(payloadSlice) != 3 {
		return nil
	}
	pendingMessageId, err := strconv.Atoi(payloadSlice[0])
	groupId, err := strconv.ParseInt(payloadSlice[1], 10, 64)
	groupTitle := payloadSlice[2]
	if err != nil {
		logger.Error("[startCaptcha] groupId err:", err)
		return nil
	}
	userId := c.Sender().ID
	pendingKey := fmt.Sprintf("%d|%d", pendingMessageId, groupId)
	record, ok := t.userCaptchaPendingTable.Get(pendingKey)
	if ok == false || record == nil || record.UserId != c.Sender().ID {
		return c.Send("您在该群没有待验证记录😁")
	}
	// 获得一个验证码
	captchaCode, imgUrl, err := captcha.GetCaptcha()
	if err != nil {
		logger.Error("[startCaptcha] get image captcha err:", err)
		return c.Send("服务器异常~，请稍后再试")
	}
	captchaMessage := fmt.Sprintf(t.settings.CaptchaMessage,
		groupTitle,
		t.settings.CaptchaTimeout,
	)
	sendMessage := &tb.Photo{
		File:    tb.FromDisk(imgUrl),
		Caption: captchaMessage,
	}
	captchaMessageMenu := &tb.ReplyMarkup{ResizeKeyboard: true}
	refreshCaptchaImageBtn := captchaMessageMenu.Data("🔁刷新验证码", "refreshCaptchaImageBtn", strconv.FormatInt(userId, 10))
	t.bot.Handle(&refreshCaptchaImageBtn, t.refreshCaptcha())
	captchaMessageMenu.Inline(
		captchaMessageMenu.Row(refreshCaptchaImageBtn),
	)
	botMsg, err := t.bot.Send(c.Chat(), sendMessage, captchaMessageMenu)
	if err != nil {
		logger.Error("[startCaptcha] send image captcha err:", err)
		return c.Send("服务器异常~，请稍后再试")
	}
	userCaptchaCodeVal := &captcha.Code{
		UserId:         userId,
		GroupId:        groupId,
		Code:           captchaCode,
		CaptchaMessage: botMsg,
		PendingMessage: record.PendingMessage,
		GroupTitle:     groupTitle,
		CreatedAt:      time.Now().Unix(),
	}
	userCaptchaCodeKey := strconv.FormatInt(userId, 10)
	t.userCaptchaCodeTable.Set(userCaptchaCodeKey, userCaptchaCodeVal)
	// TODO 需要重构，这个部分
	time.AfterFunc(time.Duration(t.settings.CaptchaTimeout)*time.Second, func() {
		_ = os.Remove(imgUrl)
		t.messageTokenMap.Remove(chatToken)
		t.userCaptchaCodeTable.Del(userCaptchaCodeKey)
		err = t.bot.Delete(botMsg)
		if err != nil {
			logger.Error("[startCaptcha] delete captcha err:", err)
		}
	})
	return nil
}

// refreshCaptcha 刷新验证码
func (t *TelegramBot) refreshCaptcha() func(c tb.Context) error {
	return func(c tb.Context) error {
		userIdStr := strconv.FormatInt(c.Sender().ID, 10)
		captchaCode, ok := t.userCaptchaCodeTable.Get(userIdStr)
		if ok == false || captchaCode == nil || captchaCode.UserId != c.Sender().ID {
			return nil
		}
		// 获得一个新验证码
		code, imgUrl, err := captcha.GetCaptcha()
		if err != nil {
			logger.Error(err)
			return c.Respond(&tb.CallbackResponse{
				Text: "服务器繁忙~",
			})
		}
		editMessage := &tb.Photo{
			File: tb.FromDisk(imgUrl),
			Caption: fmt.Sprintf(t.settings.CaptchaMessage,
				captchaCode.GroupTitle,
				t.settings.CaptchaTimeout,
			),
		}
		_, err = t.bot.Edit(c.Message(), editMessage, &tb.ReplyMarkup{InlineKeyboard: c.Message().ReplyMarkup.InlineKeyboard})
		if err != nil {
			logger.Error("[refreshCaptcha] send refreshCaptcha err:", err)
			return nil
		}
		captchaCode.Code = code
		t.userCaptchaCodeTable.Set(userIdStr, captchaCode)
		_ = os.Remove(imgUrl)
		return c.Respond(&tb.CallbackResponse{
			Text: "验证码已刷新~",
		})
	}
}

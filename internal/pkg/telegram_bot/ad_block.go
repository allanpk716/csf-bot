package telegram_bot

import (
	"fmt"
	"github.com/WQGroup/logger"
	sensitiveword "github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/sensitive_word"
	tb "gopkg.in/telebot.v3"
	"strconv"
	"strings"
)

// adBlock 广告阻止
func (t *TelegramBot) adBlock(c tb.Context) error {
	userId := c.Message().Sender.ID
	userLink := fmt.Sprintf("tg://user?id=%d", c.Message().Sender.ID)
	userNickname := c.Message().Sender.LastName + c.Message().Sender.FirstName
	messageText := c.Message().Text
	// 管理员 放行任何操作
	if t.isManage(c.Chat(), userId) {
		return nil
	}
	dict := sensitiveword.Filter.FindAll(messageText)
	if len(dict) <= 0 || len(dict) < t.settings.NumberOfForbiddenWords {
		return nil
	}
	// ban user
	restrictedUntil := t.settings.BlockTime
	if restrictedUntil <= 0 {
		restrictedUntil = tb.Forever()
	}
	err := t.bot.Restrict(c.Chat(), &tb.ChatMember{
		Rights:          tb.NoRights(),
		User:            c.Message().Sender,
		RestrictedUntil: restrictedUntil,
	})
	if err != nil {
		logger.Error("[adBlock] ban user err:", err)
		return err
	}
	// 嫌疑人窗口，可能触发的关键词
	criminalSuspectMenu := &tb.ReplyMarkup{ResizeKeyboard: true}
	blockMessage := fmt.Sprintf(t.settings.BlockMessage,
		userNickname,
		userLink,
		strings.Join(dict, ","))
	criminalSuspectBtn := criminalSuspectMenu.Data("👮🏻管理员解封", strconv.FormatInt(userId, 10))
	criminalSuspectMenu.Inline(criminalSuspectMenu.Row(criminalSuspectBtn))
	// 添加额外的广告链接
	//LoadAdMenuBtn(criminalSuspectMenu)
	t.bot.Handle(&criminalSuspectBtn, func(c tb.Context) error {
		if err = t.bot.Delete(c.Message()); err != nil {
			logger.Error("[adBlock] delete adblock message err:", err)
			return err
		}
		// 解禁用户
		err = t.bot.Restrict(c.Chat(), &tb.ChatMember{
			User:   &tb.User{ID: userId},
			Rights: tb.NoRestrictions(),
		})
		if err != nil {
			logger.Error("[adBlock] unban user err:", err)
			return err
		}
		return c.Send(fmt.Sprintf("管理员已解除对用户：[%s](%s) 的封禁", userNickname, userLink), tb.ModeMarkdownV2)
	}, t.isManageMiddleware)
	if err = c.Reply(blockMessage, criminalSuspectMenu, tb.ModeMarkdownV2); err != nil {
		logger.Error("[adBlock] reply message err:", err)
		return err
	}
	return c.Delete()
}

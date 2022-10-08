package telegram_bot

import (
	"fmt"
	"github.com/WQGroup/logger"
	tb "gopkg.in/telebot.v3"
)

func (t *TelegramBot) manulPassVerify(c tb.Context) error {

	passCode := c.Message().Payload
	// 不是私聊或者载荷为空
	if !c.Message().Private() || passCode == "" {
		return nil
	}

	_, found := t.oneTimePassCodeMap.Get(passCode)
	if found == true {

		// 移除这个一次性验证码
		t.oneTimePassCodeMap.Remove(passCode)
		// 解禁用户
		// 直接给所有群的这个用户解封
		for _, v := range t.settings.GroupIDs {
			err := t.bot.Restrict(&tb.Chat{ID: v}, &tb.ChatMember{
				User:   &tb.User{ID: c.Sender().ID},
				Rights: tb.NoRestrictions(),
			})
			if err != nil {
				logger.Error("[onTextMessage] manulPassVerify err:", err)
				return c.Send("服务器异常~，请稍后重试~")
			}
		}
		logger.Infoln("用户:", c.Sender().Username, "通过一次性验证码", passCode, "解禁")
		return c.Send("解封通过")
	} else {
		logger.Warningln("用户:", c.Sender().Username, "使用无效的一次性验证码", passCode)
		return c.Send("解封码: " + passCode + " 无效")
	}
}

func (t *TelegramBot) setOneTimePassCode(c tb.Context) error {

	passCode := c.Message().Payload
	// 不是私聊或者载荷为空
	if !c.Message().Private() || passCode == "" {
		return nil
	}

	t.oneTimePassCodeMap.Set(passCode, fmt.Sprintf("%d", c.Sender().ID))

	return c.Send(fmt.Sprintf("找这个机器人 https://t.me/%s \r\n进行对话\r\n发送下面这段代码（注意有空格）：\r\n /manual_pass_verify %s", t.bot.Me.Username, passCode))
}

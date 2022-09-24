package captcha

import (
	tb "gopkg.in/telebot.v3"
)

type Pending struct {
	UserId         int64       `json:"user_id"`
	GroupId        int64       `json:"group_id"`
	JoinAt         int64       `json:"join_at"`
	PendingMessage *tb.Message `json:"pending_message"`
}

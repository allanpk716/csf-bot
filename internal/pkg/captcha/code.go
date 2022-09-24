package captcha

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	tb "gopkg.in/telebot.v3"
)

// Code 验证码
type Code struct {
	UserId         int64       `json:"user_id"`
	GroupId        int64       `json:"group_id"`
	Code           string      `json:"code"`
	CaptchaMessage *tb.Message `json:"message_id"`
	PendingMessage *tb.Message `json:"pending_message"`
	GroupTitle     string      `json:"group_title"`
	CreatedAt      int64       `json:"created_at"`
}

type CodeTable struct {
	UserCaptchaCode cmap.ConcurrentMap[*Code]
}

func NewCodeTable() *CodeTable {
	return &CodeTable{
		UserCaptchaCode: cmap.New[*Code](),
	}
}

func (t *CodeTable) Set(key string, val *Code) {
	t.UserCaptchaCode.Set(key, val)
}

func (t *CodeTable) Get(key string) (*Code, bool) {

	return t.UserCaptchaCode.Get(key)
}

func (t *CodeTable) Del(key string) {
	t.UserCaptchaCode.Remove(key)
}

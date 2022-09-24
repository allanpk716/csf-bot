package captcha

import cmap "github.com/orcaman/concurrent-map/v2"

type PendingTable struct {
	UserCaptchaPending cmap.ConcurrentMap[*Pending]
}

func NewPendingTable() *PendingTable {
	return &PendingTable{
		UserCaptchaPending: cmap.New[*Pending](),
	}
}

func (t *PendingTable) Set(key string, val *Pending) {
	t.UserCaptchaPending.Set(key, val)
}

func (t *PendingTable) Get(key string) (*Pending, bool) {

	return t.UserCaptchaPending.Get(key)
}

func (t *PendingTable) Del(key string) {
	t.UserCaptchaPending.Remove(key)
}

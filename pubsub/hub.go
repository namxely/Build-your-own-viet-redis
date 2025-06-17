package pubsub

import (
	"github.com/namxely/Build-your-own-viet-redis/datastruct/dict"
	"github.com/namxely/Build-your-own-viet-redis/datastruct/lock"
)

// Hub stores all subscribe relations
type Hub struct {
	// channel -> list(*Client)
	subs dict.Dict
	// lock channel
	subsLocker *lock.Locks
}

// MakeHub creates new hub
func MakeHub() *Hub {
	return &Hub{
		subs:       dict.MakeConcurrent(4),
		subsLocker: lock.Make(16),
	}
}

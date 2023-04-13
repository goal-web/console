package scheduling

import (
	"github.com/goal-web/contracts"
)

type Mutex struct {
	redis contracts.RedisConnection
}

func NewMutex(connection contracts.RedisConnection) *Mutex {
	return &Mutex{redis: connection}
}

func (mutex *Mutex) Create(event *Event) bool {
	if mutex.redis != nil {
		var result, err = mutex.redis.SetNX(event.MutexName(), "1", event.expiresAt)
		return err == nil && result
	}
	return true
}
func (mutex *Mutex) Exists(event *Event) bool {
	if mutex.redis != nil {
		var num, err = mutex.redis.Exists(event.MutexName())
		return err == nil && num > 0
	}
	return true
}

func (mutex *Mutex) Forget(event *Event) {
	if mutex.redis != nil {
		_, _ = mutex.redis.Del(event.MutexName())
	}
}

package scheduling

import (
	"github.com/goal-web/contracts"
)

type Mutex struct {
	redis contracts.RedisFactory
	store string
}

func (mutex *Mutex) Create(event *Event) bool {
	if mutex.redis != nil {
		var _, err = mutex.redis.Connection(mutex.store).Set(event.MutexName(), "1", event.expiresAt)
		return err == nil
	}
	return true
}
func (mutex *Mutex) Exists(event *Event) bool {
	if mutex.redis != nil {
		var num, err = mutex.redis.Connection(mutex.store).Exists(event.MutexName())
		return err == nil && num > 0
	}
	return true
}

func (mutex *Mutex) Forget(event *Event) {
	if mutex.redis != nil {
		_, _ = mutex.redis.Connection(mutex.store).Del(event.MutexName())
	}
}

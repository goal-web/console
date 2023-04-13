package scheduling

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func NewCallbackEvent(mutex *Mutex, callback any, timezone string) contracts.CallbackEvent {
	return &CallbackEvent{
		Event:       NewEvent(mutex, callback, timezone),
		description: "",
	}
}

type CallbackEvent struct {
	*Event
	description string
}

func (event *CallbackEvent) Description(description string) contracts.CallbackEvent {
	event.description = description
	return event
}

func (event *CallbackEvent) MutexName() string {
	if event.mutexName == "" {
		return fmt.Sprintf("goal.schedule-%s", utils.Md5(event.expression+event.description))
	}
	return event.mutexName
}

package scheduling

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func NewCommandEvent(command string, mutex *Mutex, callback any, timezone string) contracts.CommandEvent {
	return &CommandEvent{
		Event:   NewEvent(mutex, callback, timezone),
		command: command,
	}
}

type CommandEvent struct {
	*Event
	command string
}

func (event *CommandEvent) Command(command string) contracts.CommandEvent {
	event.command = command
	return event
}

func (event *CommandEvent) MutexName() string {
	if event.mutexName == "" {
		return fmt.Sprintf("goal.schedule-%s", utils.Md5(event.expression+event.command))
	}
	return event.mutexName
}

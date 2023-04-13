package scheduling

import (
	"github.com/goal-web/application"
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"os/exec"
)

type Schedule struct {
	store    string
	timezone string
	mutex    *Mutex
	app      contracts.Application

	events []contracts.ScheduleEvent
}

func (schedule *Schedule) GetEvents() []contracts.ScheduleEvent {
	return schedule.events
}

func (schedule *Schedule) UseStore(store string) {
	schedule.store = store
}

func NewSchedule(app contracts.Application) contracts.Schedule {
	var (
		appConfig = app.Get("config").(contracts.Config).Get("app").(application.Config)
		redis, _  = app.Get("redis.factory").(contracts.RedisFactory)
	)
	return &Schedule{
		timezone: appConfig.Timezone,
		mutex:    &Mutex{redis: redis.Connection()},
		app:      app,
		events:   make([]contracts.ScheduleEvent, 0),
	}
}

func (schedule *Schedule) Call(callback any, args ...any) contracts.CallbackEvent {
	event := NewCallbackEvent(schedule.mutex, func() {
		schedule.app.Call(callback, args...)
	}, schedule.timezone)
	schedule.events = append(schedule.events, event)
	return event
}

func (schedule *Schedule) Command(command contracts.Command, args ...string) contracts.CommandEvent {
	args = append([]string{command.GetName()}, args...)
	input := inputs.String(args...)
	err := command.InjectArguments(input.GetArguments())
	if err != nil {
		logs.WithError(err).WithField("args", args).Debug("Schedule.Command: arguments invalid")
		panic(err) // 因为这个阶段框架还没正式运行，所以 panic
	}
	event := NewCommandEvent(command.GetName(), schedule.mutex, func(console contracts.Console) {
		command.Handle()
	}, schedule.timezone)
	schedule.events = append(schedule.events, event)
	return event
}

func (schedule *Schedule) Exec(command string, args ...string) contracts.CommandEvent {
	var event = NewCommandEvent(command, schedule.mutex, func(console contracts.Console) {
		if console.Exists(command) {
			args = append([]string{command}, args...)
			input := inputs.String(args...)
			console.Run(&input)
		} else {
			if err := exec.Command(command, args...).Run(); err != nil {
				logs.WithError(err).
					WithField("command", command).
					WithField("args", args).
					Error("Schedule.Exec: failed")
			}
		}

	}, schedule.timezone)
	schedule.events = append(schedule.events, event)
	return event
}

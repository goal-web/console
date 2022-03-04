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

func (this *Schedule) GetEvents() []contracts.ScheduleEvent {
	return this.events
}

func (this *Schedule) UseStore(store string) {
	this.store = store
}

func NewSchedule(app contracts.Application) contracts.Schedule {
	var (
		appConfig = app.Get("config").(contracts.Config).Get("app").(application.Config)
		redis, _  = app.Get("redis.factory").(contracts.RedisFactory)
	)
	return &Schedule{
		timezone: appConfig.Timezone,
		mutex: &Mutex{
			redis: redis,
			store: "cache",
		},
		app:    app,
		events: make([]contracts.ScheduleEvent, 0),
	}
}

func (this *Schedule) Call(callback interface{}, args ...interface{}) contracts.CallbackEvent {
	event := NewCallbackEvent(this.mutex, func() {
		this.app.Call(callback, args...)
	}, this.timezone)
	this.events = append(this.events, event)
	return event
}

func (this *Schedule) Command(command contracts.Command, args ...string) contracts.CommandEvent {
	args = append([]string{command.GetName()}, args...)
	input := inputs.StringArray(args)
	err := command.InjectArguments(input.GetArguments())
	if err != nil {
		logs.WithError(err).WithField("args", args).Debug("Schedule.Command: arguments invalid")
		panic(err) // 因为这个阶段框架还没正式运行，所以 panic
	}
	event := NewCommandEvent(command.GetName(), this.mutex, func(console contracts.Console) {
		command.Handle()
	}, this.timezone)
	this.events = append(this.events, event)
	return event
}

func (this *Schedule) Exec(command string, args ...string) contracts.CommandEvent {
	var event = NewCommandEvent(command, this.mutex, func(console contracts.Console) {
		if console.Exists(command) {
			args = append([]string{command}, args...)
			input := inputs.StringArray(args)
			console.Run(&input)
		} else {
			if err := exec.Command(command, args...).Run(); err != nil {
				logs.WithError(err).
					WithField("command", command).
					WithField("args", args).
					Debug("Schedule.Exec: failed")
			}
		}

	}, this.timezone)
	this.events = append(this.events, event)
	return event
}

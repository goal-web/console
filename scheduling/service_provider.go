package scheduling

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/gorhill/cronexpr"
	"time"
)

type serviceProvider struct {
	intervalCloseChannel chan bool
	runningErrChannel    chan error
	app                  contracts.Application
	execRecords          map[int]time.Time
	exceptionHandler     contracts.ExceptionHandler
}

func NewService() contracts.ServiceProvider {
	return &serviceProvider{}
}

func (provider *serviceProvider) Register(application contracts.Application) {
	provider.app = application
}

func (provider *serviceProvider) runScheduleEvents(events []contracts.ScheduleEvent) {
	if len(events) > 0 {
		// 并发执行所有事件
		now := time.Now()
		for index, event := range events {
			lastExecTime := provider.execRecords[index]
			nextTime := cronexpr.MustParse(event.Expression()).Next(lastExecTime)
			diff := now.Sub(nextTime).Seconds()
			if diff >= 0 && int(diff) == 0 {
				provider.execRecords[index] = now
				go (func(event contracts.ScheduleEvent) {
					defer func() {
						if err := recover(); err != nil {
							provider.exceptionHandler.Handle(&ScheduleEventException{
								Err:      errors.New("task execution failed"),
								Previous: exceptions.WithRecover(err),
							})
						}
					}()
					event.Run(provider.app)
				})(event)
			} else if nextTime.Before(now) {
				provider.execRecords[index] = now
			}
		}
	}
}

func (provider *serviceProvider) Start() error {
	provider.execRecords = make(map[int]time.Time)
	provider.app.Call(func(schedule contracts.Schedule, exceptionHandler contracts.ExceptionHandler) {
		provider.exceptionHandler = exceptionHandler
		if len(schedule.GetEvents()) > 0 {
			provider.runningErrChannel = make(chan error)
			provider.intervalCloseChannel = utils.SetInterval(1, func() {
				provider.runScheduleEvents(schedule.GetEvents())
			}, func() {
				logs.Default().Info("the goal scheduling is closed")
			})
		}
	})
	return <-provider.runningErrChannel
}

func (provider *serviceProvider) Stop() {
	if provider.intervalCloseChannel != nil {
		provider.intervalCloseChannel <- true
		close(provider.intervalCloseChannel)
	}
	if provider.runningErrChannel != nil {
		provider.runningErrChannel <- nil
		close(provider.runningErrChannel)
	}
}

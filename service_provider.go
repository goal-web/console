package console

import (
	"github.com/goal-web/application"
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/golang-module/carbon/v2"
	"github.com/gorhill/cronexpr"
	"reflect"
	"time"
)

type Provider func(application contracts.Application) contracts.Console

type serviceProvider struct {
	ConsoleProvider Provider

	stopChan         chan bool
	serverIdChan     chan bool
	app              contracts.Application
	execRecords      map[int]time.Time
	exceptionHandler contracts.ExceptionHandler
}

func NewService(provider Provider) contracts.ServiceProvider {
	return &serviceProvider{ConsoleProvider: provider}
}

func (provider *serviceProvider) Register(application contracts.Application) {
	provider.app = application
	provider.exceptionHandler = application.Get("exceptions.handler").(contracts.ExceptionHandler)

	application.Singleton("console", func() contracts.Console {
		console := provider.ConsoleProvider(application)
		console.Schedule(console.GetSchedule())
		return console
	})
	application.Singleton("scheduling", func(console contracts.Console) contracts.Schedule {
		return console.GetSchedule()
	})
	application.Singleton("console.input", func() contracts.ConsoleInput {
		return inputs.NewOSArgsInput()
	})
}

func (provider *serviceProvider) runScheduleEvents(events []contracts.ScheduleEvent) {
	if len(events) > 0 {
		// 并发执行所有事件
		now := time.Now()
		for index, event := range events {
			lastExecTime := provider.execRecords[index]
			nextTime := carbon.Time2Carbon(cronexpr.MustParse(event.Expression()).Next(lastExecTime))
			nowCarbon := carbon.Time2Carbon(now)
			if nextTime.DiffInSeconds(nowCarbon) == 0 {
				provider.execRecords[index] = now
				go (func(event contracts.ScheduleEvent) {
					defer func() {
						if err := recover(); err != nil {
							provider.exceptionHandler.Handle(ScheduleEventException{
								Exception: exceptions.WithRecover(err, contracts.Fields{
									"expression": event.Expression(),
									"mutex_name": event.MutexName(),
									"one_server": event.OnOneServer(),
									"event":      utils.GetTypeKey(reflect.TypeOf(event)),
								}),
							})
						}
					}()
					event.Run(provider.app)
				})(event)
			} else if nextTime.Lt(nowCarbon) {
				provider.execRecords[index] = now
			}
		}
	}
}

func (provider *serviceProvider) Start() error {
	provider.execRecords = make(map[int]time.Time)
	go provider.maintainServerId()
	provider.app.Call(func(schedule contracts.Schedule) {
		if len(schedule.GetEvents()) > 0 {
			provider.stopChan = utils.SetInterval(1, func() {
				provider.runScheduleEvents(schedule.GetEvents())
			}, func() {
				logs.Default().Info("the goal scheduling is closed")
			})
		}
	})
	return nil
}

func (provider *serviceProvider) Stop() {
	if provider.stopChan != nil {
		provider.stopChan <- true
	}
	if provider.serverIdChan != nil {
		provider.serverIdChan <- true
	}
}

// maintainServerId 维护服务实例ID
func (provider *serviceProvider) maintainServerId() {
	provider.app.Call(func(redis contracts.RedisConnection, config contracts.Config, handler contracts.ExceptionHandler) {
		appConfig := config.Get("app").(application.Config)
		provider.serverIdChan = utils.SetInterval(1, func() {
			// 维持当前服务心跳
			_, _ = redis.Set("goal.server."+appConfig.ServerId, time.Now().String(), time.Second*2)
		}, func() {
			_, _ = redis.Del("goal.server." + appConfig.ServerId)
		})
	})
}

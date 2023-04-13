package console

import (
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
)

type Provider func(application contracts.Application) contracts.Console

type serviceProvider struct {
	ConsoleProvider Provider
	app             contracts.Application
}

func NewService(provider Provider) contracts.ServiceProvider {
	return &serviceProvider{ConsoleProvider: provider}
}

func (provider *serviceProvider) Register(application contracts.Application) {
	provider.app = application

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

func (provider *serviceProvider) Start() error {
	return nil
}

func (provider *serviceProvider) Stop() {
}

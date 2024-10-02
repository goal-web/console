package console

import (
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
)

type serviceProvider struct {
	app contracts.Application
}

func NewService() contracts.ServiceProvider {
	return &serviceProvider{}
}

func (provider *serviceProvider) Register(application contracts.Application) {
	provider.app = application

	application.Singleton("console", func() contracts.Console {
		return NewKernel(application, nil)
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

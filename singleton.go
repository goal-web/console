package console

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

var singleton contracts.Console
var once sync.Once

func Default() contracts.Console {
	once.Do(func() {
		singleton = application.Get("console").(contracts.Console)
	})

	return singleton
}

func Call(command string, arguments contracts.CommandArguments) any {
	return Default().Call(command, arguments)
}

func Run(input contracts.ConsoleInput) any {
	return Default().Run(input)
}

func Exists(name string) bool {
	return Default().Exists(name)
}

func RegisterCommand(command contracts.CommandProvider) {
	Default().RegisterCommand(command)
}

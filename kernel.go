package console

import (
	"errors"
	"fmt"
	"github.com/goal-web/console/scheduling"
	"github.com/goal-web/contracts"
	"github.com/modood/table"
)

var CommandDontExists = errors.New("命令不存在！")

const logoText = "  ▄████  ▒█████   ▄▄▄       ██▓    \n ██▒ ▀█▒▒██▒  ██▒▒████▄    ▓██▒    \n▒██░▄▄▄░▒██░  ██▒▒██  ▀█▄  ▒██░    \n░▓█  ██▓▒██   ██░░██▄▄▄▄██ ▒██░    \n░▒▓███▀▒░ ████▓▒░ ▓█   ▓██▒░██████▒\n ░▒   ▒ ░ ▒░▒░▒░  ▒▒   ▓▒█░░ ▒░▓  ░\n  ░   ░   ░ ▒ ▒░   ▒   ▒▒ ░░ ░ ▒  ░\n░ ░   ░ ░ ░ ░ ▒    ░   ▒     ░ ░   \n      ░     ░ ░        ░  ░    ░  ░\n                                   "

type Kernel struct {
	app              contracts.Application
	commands         map[string]contracts.CommandProvider
	schedule         contracts.Schedule
	exceptionHandler contracts.ExceptionHandler
}

func NewKernel(app contracts.Application, commandProviders []contracts.CommandProvider) *Kernel {
	var commands = make(map[string]contracts.CommandProvider)
	for _, provider := range commandProviders {
		commands[provider(app).GetName()] = provider
	}
	return &Kernel{
		app:              app,
		commands:         commands,
		schedule:         scheduling.NewSchedule(app),
		exceptionHandler: app.Get("exceptions.handler").(contracts.ExceptionHandler),
	}
}

func (kernel *Kernel) RegisterCommand(name string, command contracts.CommandProvider) {
	kernel.commands[name] = command
}

func (kernel *Kernel) GetSchedule() contracts.Schedule {
	return kernel.schedule
}

func (kernel *Kernel) Schedule(schedule contracts.Schedule) {
}

type CommandItem struct {
	Command     string
	Signature   string
	Description string
}

func (kernel *Kernel) Help() {
	cmdTable := make([]CommandItem, 0)
	for _, command := range kernel.commands {
		cmd := command(kernel.app)
		cmdTable = append(cmdTable, CommandItem{
			Command:     cmd.GetName(),
			Signature:   cmd.GetSignature(),
			Description: cmd.GetDescription(),
		})
	}
	fmt.Println(logoText)
	table.Output(cmdTable)
}

func (kernel *Kernel) Call(cmd string, arguments contracts.CommandArguments) any {
	if cmd == "" {
		kernel.Help()
		return nil
	}
	for name, provider := range kernel.commands {
		if cmd == name {
			command := provider(kernel.app)
			if arguments.Exists("h") || arguments.Exists("help") {
				fmt.Println(logoText)
				fmt.Printf(" %s 命令：%s\n", command.GetName(), command.GetDescription())
				fmt.Println(command.GetHelp())
				return nil
			}
			if err := command.InjectArguments(arguments); err != nil {
				kernel.exceptionHandler.Handle(&CommandArgumentException{Err: errors.New("the parameter is wrong")})
				fmt.Println(err.Error())
				fmt.Println(command.GetHelp())
				return nil
			}
			return command.Handle()
		}
	}
	return CommandDontExists
}

func (kernel *Kernel) Run(input contracts.ConsoleInput) any {
	return kernel.Call(input.GetCommand(), input.GetArguments())
}

func (kernel *Kernel) Exists(name string) bool {
	return kernel.commands[name] != nil
}

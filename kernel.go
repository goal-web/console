package console

import (
	"errors"
	"fmt"
	"github.com/goal-web/console/scheduling"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
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

func (this *Kernel) RegisterCommand(name string, command contracts.CommandProvider) {
	this.commands[name] = command
}

func (this *Kernel) GetSchedule() contracts.Schedule {
	return this.schedule
}

func (this *Kernel) Schedule(schedule contracts.Schedule) {
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

type CommandItem struct {
	Command     string
	Signature   string
	Description string
}

func (this Kernel) Help() {
	cmdTable := make([]CommandItem, 0)
	for _, command := range this.commands {
		cmd := command(this.app)
		cmdTable = append(cmdTable, CommandItem{
			Command:     cmd.GetName(),
			Signature:   cmd.GetSignature(),
			Description: cmd.GetDescription(),
		})
	}
	fmt.Println(logoText)
	table.Output(cmdTable)
}

func (this *Kernel) Call(cmd string, arguments contracts.CommandArguments) interface{} {
	if cmd == "" {
		this.Help()
		return nil
	}
	for name, provider := range this.commands {
		if cmd == name {
			command := provider(this.app)
			if arguments.Exists("h") || arguments.Exists("help") {
				fmt.Println(logoText)
				fmt.Printf(" %s 命令：%s\n", command.GetName(), command.GetDescription())
				fmt.Println(command.GetHelp())
				return nil
			}
			if err := command.InjectArguments(arguments); err != nil {
				this.exceptionHandler.Handle(CommandArgumentException{
					exceptions.WithError(err, contracts.Fields{
						"command":   cmd,
						"arguments": arguments,
					}),
				})
				fmt.Println(err.Error())
				fmt.Println(command.GetHelp())
				return nil
			}
			return command.Handle()
		}
	}
	return CommandDontExists
}

func (this *Kernel) Run(input contracts.ConsoleInput) interface{} {
	return this.Call(input.GetCommand(), input.GetArguments())
}

func (this *Kernel) Exists(name string) bool {
	return this.commands[name] != nil
}

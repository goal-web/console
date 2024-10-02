package console

import (
	"errors"
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/modood/table"
)

var CommandDontExists = errors.New("命令不存在！")

const logoText = "  ▄████  ▒█████   ▄▄▄       ██▓    \n ██▒ ▀█▒▒██▒  ██▒▒████▄    ▓██▒    \n▒██░▄▄▄░▒██░  ██▒▒██  ▀█▄  ▒██░    \n░▓█  ██▓▒██   ██░░██▄▄▄▄██ ▒██░    \n░▒▓███▀▒░ ████▓▒░ ▓█   ▓██▒░██████▒\n ░▒   ▒ ░ ▒░▒░▒░  ▒▒   ▓▒█░░ ▒░▓  ░\n  ░   ░   ░ ▒ ▒░   ▒   ▒▒ ░░ ░ ▒  ░\n░ ░   ░ ░ ░ ░ ▒    ░   ▒     ░ ░   \n      ░     ░ ░        ░  ░    ░  ░\n                                   "

type Kernel struct {
	app              contracts.Application
	commands         map[string]contracts.Command
	handlers         map[string]contracts.CommandHandlerProvider
	exceptionHandler contracts.ExceptionHandler
}

func NewKernel(app contracts.Application, commandProviders []contracts.CommandProvider) *Kernel {
	var handlers = make(map[string]contracts.CommandHandlerProvider)
	var commands = make(map[string]contracts.Command)
	for _, provider := range commandProviders {
		cmd, handlerProvider := provider()
		handlers[cmd.GetName()] = handlerProvider
		commands[cmd.GetName()] = cmd
	}
	return &Kernel{
		app:              app,
		handlers:         handlers,
		commands:         commands,
		exceptionHandler: app.Get("exceptions.handler").(contracts.ExceptionHandler),
	}
}

func (kernel *Kernel) RegisterCommand(command contracts.CommandProvider) {
	cmd, handlerProvider := command()
	kernel.handlers[cmd.GetName()] = handlerProvider
	kernel.commands[cmd.GetName()] = cmd
}

type CommandItem struct {
	Command     string
	Signature   string
	Description string
}

func (kernel *Kernel) Help() {
	cmdTable := make([]CommandItem, 0)
	for _, cmd := range kernel.commands {
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
	command, ok := kernel.commands[cmd]
	if !ok {
		return CommandDontExists
	}
	if arguments.Exists("h") || arguments.Exists("help") {
		fmt.Println(logoText)
		fmt.Printf(" %s 命令：%s\n", command.GetName(), command.GetDescription())
		fmt.Println(command.GetHelp())
		return nil
	}
	handler := kernel.handlers[cmd](kernel.app)
	if err := handler.InjectArguments(command.GetArgs(), arguments); err != nil {
		kernel.exceptionHandler.Handle(&CommandArgumentException{Err: errors.New("the parameter is wrong")})
		fmt.Println(err.Error())
		fmt.Println(command.GetHelp())
		return nil
	}
	return handler.Handle()
}

func (kernel *Kernel) Run(input contracts.ConsoleInput) any {
	return kernel.Call(input.GetCommand(), input.GetArguments())
}

func (kernel *Kernel) Exists(name string) bool {
	return kernel.commands[name] != nil
}

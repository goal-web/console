package inputs

import (
	"github.com/goal-web/console/arguments"
	"github.com/goal-web/contracts"
	"strings"
)

type StringArrayInput struct {
	ArgsArray []string
}

func StringArray(argsArray []string) StringArrayInput {
	return StringArrayInput{argsArray}
}

func (input *StringArrayInput) GetCommand() string {
	if len(input.ArgsArray) > 0 {
		return input.ArgsArray[0]
	}
	return ""
}

func (input *StringArrayInput) GetArguments() contracts.CommandArguments {
	if len(input.ArgsArray) > 0 {
		args := make([]string, 0)
		options := contracts.Fields{}

		for _, arg := range input.ArgsArray[1:] {
			if strings.HasPrefix(arg, "--") {
				if argArr := strings.Split(strings.ReplaceAll(arg, "--", ""), "="); len(argArr) > 1 {
					options[argArr[0]] = argArr[1]
				} else {
					options[argArr[0]] = true
				}
			} else if strings.HasPrefix(arg, "-") {
				if argArr := strings.Split(strings.ReplaceAll(arg, "-", ""), "="); len(argArr) > 1 {
					options[argArr[0]] = argArr[1]
				} else {
					options[argArr[0]] = true
				}
			} else {
				args = append(args, arg)
			}
		}

		return arguments.NewArguments(args, options)
	}
	return arguments.NewArguments([]string{}, nil)
}

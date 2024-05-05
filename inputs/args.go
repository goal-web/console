package inputs

import (
	"github.com/goal-web/contracts"
	"os"
	"strings"
)

type ArgsInput struct {
	StringArrayInput
}

func NewOSArgsInput() contracts.ConsoleInput {
	var args []string
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && len(args) == 0 {
			continue
		}
		args = append(args, arg)
	}
	return &ArgsInput{String(args...)}
}

package arguments

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"strings"
)

type Arguments struct {
	supports.BaseFields
	args    []string
	options contracts.Fields
}

func NewArguments(args []string, options contracts.Fields) contracts.CommandArguments {
	arguments := &Arguments{
		args:       args,
		BaseFields: supports.BaseFields{},
		options:    options,
	}

	arguments.BaseFields.Provider = arguments
	return arguments
}

func (args *Arguments) GetArg(index int) string {
	if index >= len(args.args) {
		return ""
	}
	return args.args[index]
}

func (args *Arguments) GetArgs() []string {
	return args.args
}
func (args *Arguments) SetOption(key string, value any) {
	args.ToFields()[key] = value
}

func (args *Arguments) StringArrayOption(key string, defaultValue []string) []string {
	if value := args.GetString(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func (args *Arguments) Int64ArrayOption(key string, defaultValue []int64) []int64 {
	if value := args.GetString(key); value != "" {
		values := make([]int64, 0)
		for _, value = range strings.Split(value, ",") {
			values = append(values, utils.ToInt64(value, 0))
		}
		return values
	}
	return defaultValue
}

func (args *Arguments) IntArrayOption(key string, defaultValue []int) []int {
	if value := args.GetString(key); value != "" {
		values := make([]int, 0)
		for _, value = range strings.Split(value, ",") {
			values = append(values, utils.ToInt(value, 0))
		}
		return values
	}
	return defaultValue
}

func (args *Arguments) Float64ArrayOption(key string, defaultValue []float64) []float64 {
	if value := args.GetString(key); value != "" {
		values := make([]float64, 0)
		for _, value = range strings.Split(value, ",") {
			values = append(values, utils.ToFloat64(value, 0))
		}
		return values
	}
	return defaultValue
}

func (args *Arguments) FloatArrayOption(key string, defaultValue []float32) []float32 {
	if value := args.GetString(key); value != "" {
		values := make([]float32, 0)
		for _, value = range strings.Split(value, ",") {
			values = append(values, utils.ToFloat(value, 0))
		}
		return values
	}
	return defaultValue
}

func (args *Arguments) ToFields() contracts.Fields {
	return args.options
}

func (args *Arguments) Exists(key string) bool {
	_, exists := args.options[key]
	return exists
}

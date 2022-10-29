package injection

import "go.uber.org/fx"

var options []fx.Option

func Provide(constructors ...interface{}) {
	options = append(options, fx.Provide(constructors...))
}

func Invoke(funcs ...interface{}) {
	options = append(options, fx.Invoke(funcs...))
}

func GetOptions() []fx.Option {
	return options
}

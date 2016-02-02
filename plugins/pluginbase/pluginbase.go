package pluginbase

import "github.com/robertkrimen/otto"

type Plugin struct {
	Name     string
	Init     PluginInit
	Disabled bool
}
type Result struct {
	Suc interface{}
	Err interface{}
}

type JSCall func(otto.FunctionCall) otto.Value
type PluginInit func(*otto.Otto)

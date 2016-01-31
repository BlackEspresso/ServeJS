package pluginbase

import "github.com/robertkrimen/otto"

type Plugin struct {
	Name     string
	Init     PluginInit
	Disabled bool
}

type JSCall func(otto.FunctionCall) otto.Value
type PluginInit func(*otto.Otto)

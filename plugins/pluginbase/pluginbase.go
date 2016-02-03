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

func ToResult(vm *otto.Otto, success interface{}, err error) otto.Value {
	res := Result{}
	if err != nil {
		res.Err = err.Error()
	} else {
		res.Suc = success
	}
	resV, _ := vm.ToValue(res)
	return resV
}

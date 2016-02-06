package pluginbase

import (
	"net/http"

	"github.com/robertkrimen/otto"
)

type FuncMapping map[string]func(w http.ResponseWriter, r *http.Request)

type Plugin struct {
	Name        string
	Init        PluginInit
	Disabled    bool
	HttpMapping map[string]func(w http.ResponseWriter, r *http.Request)
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

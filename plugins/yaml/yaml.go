package yaml

import (
	"./../pluginbase"
	"github.com/robertkrimen/otto"
	"gopkg.in/gomail.v1"
)

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "yaml",
		Init: func(vm *otto.Otto) {

			vm.Set("send", func(c otto.FunctionCall) otto.Value {

				return otto.TrueValue()
			})
		},
	}

	return &p1
}

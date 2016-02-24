package time

import (
	"time"

	"./../modules"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "time",
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *modules.JsVm) otto.Value {
	obj, _ := vm.Object("({})")
	obj.Set("sleep", func(c otto.FunctionCall) otto.Value {
		sec, _ := c.Argument(0).ToFloat()
		time.Sleep(time.Duration(sec) * time.Millisecond)
		return otto.TrueValue()
	})
	return obj.Value()
}

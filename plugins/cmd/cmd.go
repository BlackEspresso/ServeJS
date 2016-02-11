package cmd

import (
	"bytes"
	"os/exec"

	"./../modules"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "cmd",
		Init: func(vm *otto.Otto) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("runCmd", func(c otto.FunctionCall) otto.Value {
				list := []string{}
				for _, v := range c.ArgumentList {
					str, _ := v.ToString()
					list = append(list, str)
				}

				cmd := exec.Command(list[0], list[1:]...)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				return modules.ToResult(vm, out.String(), err)
			})
			return obj.Value()
		},
	}

	return &p1
}

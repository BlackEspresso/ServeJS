package cmd

import (
	"bytes"
	"os/exec"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "cmd",
		Init: func(vm *otto.Otto) {
			vm.Set("runCmd", func(c otto.FunctionCall) otto.Value {
				list := []string{}
				for _, v := range c.ArgumentList {
					str, _ := v.ToString()
					list = append(list, str)
				}

				cmd := exec.Command(list[0], list[1:]...)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				return pluginbase.ToResult(vm, out.String(), err)
			})
		},
	}

	return &p1
}

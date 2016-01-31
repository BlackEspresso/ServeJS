package file

import (
	"io/ioutil"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "file",
		Init: func(vm *otto.Otto) {
			vm.Set("writeFile", writeFile)
			vm.Set("readFile", writeFile)
		},
	}
	return &p
}

func writeFile(c otto.FunctionCall) otto.Value {
	folder, _ := c.Argument(0).ToString()
	file, _ := c.Argument(1).ToString()
	data, _ := c.Argument(2).ToString()
	err := ioutil.WriteFile("./"+folder+"/"+file, []byte(data), 755)
	if err != nil {
		k, _ := otto.ToValue(err.Error())
		return k
	}
	return otto.UndefinedValue()
}

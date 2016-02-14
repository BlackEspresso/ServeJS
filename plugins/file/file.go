package file

import (
	"io/ioutil"

	"./../modules"
	"github.com/robertkrimen/otto"
)

type JsFileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

func InitPlugin() *modules.Plugin {
	p := modules.Plugin{
		Name: "file",
		Init: func(vm *otto.Otto) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("writeFile", writeFile)

			obj.Set("readFile", func(c otto.FunctionCall) otto.Value {
				path, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				data, err := ioutil.ReadFile(path)

				return modules.ToResult(vm, string(data), err)
			})

			obj.Set("readDir", func(c otto.FunctionCall) otto.Value {
				folder, _ := c.Argument(0).ToString()
				fileInfos, err := ioutil.ReadDir("./" + folder)

				if err != nil {
					return modules.ToResult(vm, nil, err)
				}

				jsFileInfos := []*JsFileInfo{}

				for _, v := range fileInfos {
					fi := JsFileInfo{
						v.Name(),
						v.Size(),
						v.IsDir(),
					}
					jsFileInfos = append(jsFileInfos, &fi)
				}

				return modules.ToResult(vm, jsFileInfos, err)
			})
			return obj.Value()
		},
	}
	return &p
}

func writeFile(c otto.FunctionCall) otto.Value {
	path, err := c.Argument(0).ToString()
	if err != nil {
		return modules.ToResult(c.Otto, nil, err)
	}

	data, _ := c.Argument(1).ToString()
	//perm, _ := c.Argument(3).ToString()
	err = ioutil.WriteFile(path, []byte(data), 777)
	return modules.ToResult(c.Otto, nil, err)
}

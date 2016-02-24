package file

import (
	"io/ioutil"
	"os"

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
		Init: func(vm *modules.JsVm) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("write", func(c otto.FunctionCall) otto.Value {
				path, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}

				data, _ := c.Argument(1).ToString()
				//perm, _ := c.Argument(3).ToString()
				err = ioutil.WriteFile(path, []byte(data), 777)
				return modules.ToResult(vm, nil, err)
			})

			obj.Set("read", func(c otto.FunctionCall) otto.Value {
				path, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				data, err := ioutil.ReadFile(path)

				return modules.ToResult(vm, string(data), err)
			})

			obj.Set("move", func(c otto.FunctionCall) otto.Value {
				source, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				target, err := c.Argument(1).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				err = os.Rename(source, target)
				return modules.ToResult(vm, true, err)
			})

			obj.Set("remove", func(c otto.FunctionCall) otto.Value {
				source, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}

				err = os.Remove(source)
				return modules.ToResult(vm, true, err)
			})

			obj.Set("removeAll", func(c otto.FunctionCall) otto.Value {
				source, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}

				err = os.RemoveAll(source)
				return modules.ToResult(vm, true, err)
			})

			obj.Set("mkdirAll", func(c otto.FunctionCall) otto.Value {
				path, err := c.Argument(0).ToString()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				perm, err := c.Argument(1).ToInteger()
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}

				err = os.MkdirAll(path, os.FileMode(perm))
				return modules.ToResult(vm, true, err)
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

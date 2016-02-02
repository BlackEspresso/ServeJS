package file

import (
	"io/ioutil"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

type JsFileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

func InitPlugin() *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "file",
		Init: func(vm *otto.Otto) {
			vm.Set("writeFile", writeFile)

			vm.Set("readFile", func(c otto.FunctionCall) otto.Value {
				folder, _ := c.Argument(0).ToString()
				file, _ := c.Argument(1).ToString()
				data, err := ioutil.ReadFile("./" + folder + "/" + file)
				res := pluginbase.Result{}
				if err != nil {
					res.Err = err.Error()
				} else {
					res.Suc = string(data)
				}
				resV, _ := vm.ToValue(res)
				return resV
			})

			vm.Set("readDir", func(c otto.FunctionCall) otto.Value {
				folder, _ := c.Argument(0).ToString()
				fileInfos, err := ioutil.ReadDir("./" + folder)
				res := pluginbase.Result{}

				if err != nil {
					res.Err = err.Error()
				} else {

					jsFileInfos := []*JsFileInfo{}

					for _, v := range fileInfos {
						fi := JsFileInfo{
							v.Name(),
							v.Size(),
							v.IsDir(),
						}
						jsFileInfos = append(jsFileInfos, &fi)
					}
					res.Suc = jsFileInfos

				}
				resV, _ := vm.ToValue(res)
				return resV
			})
		},
	}
	return &p
}

func writeFile(c otto.FunctionCall) otto.Value {
	folder, _ := c.Argument(0).ToString()
	file, _ := c.Argument(1).ToString()
	data, _ := c.Argument(2).ToString()
	//perm, _ := c.Argument(3).ToString()
	err := ioutil.WriteFile("./"+folder+"/"+file, []byte(data), 777)
	if err != nil {
		k, _ := otto.ToValue(err.Error())
		return k
	}
	return otto.UndefinedValue()
}

package settings

import (
	"encoding/json"
	"io/ioutil"

	"./../cache"
	"./../modules"
	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port    int
	Plugins map[string]map[string]string
}

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "settings",
		Init: func(vm *modules.JsVm) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("read", func(c otto.FunctionCall) otto.Value {
				path, _ := c.Argument(0).ToString()
				configraw, err := ioutil.ReadFile(path)
				if err != nil {
					return modules.ToResult(vm, nil, err)
				}
				var serverConf = Configuration{}
				err = yaml.Unmarshal(configraw, &serverConf)
				return modules.ToResult(vm, serverConf, err)
			})
			return obj.Value()
		},
	}

	return &p1
}

func GetSettings() Configuration {
	k := cache.GetCache()["settings"]
	var settings = Configuration{}
	json.Unmarshal([]byte(k), &settings)
	return settings
}

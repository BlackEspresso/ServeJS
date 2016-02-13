package cache

import (
	"encoding/json"
	"io/ioutil"

	"./../modules"
	"github.com/robertkrimen/otto"
)

var kvCache map[string]string = map[string]string{}

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "cache",
		Init: registerVM,
	}

	return &p1
}

func GetCache() map[string]string {
	return kvCache
}

func registerVM(vm *otto.Otto) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("set", func(c otto.FunctionCall) otto.Value {
		key, _ := c.Argument(0).ToString()
		val, _ := c.Argument(1).ToString()
		kvCache[key] = val
		return otto.TrueValue()
	})
	obj.Set("get", func(c otto.FunctionCall) otto.Value {
		key, _ := c.Argument(0).ToString()
		val, ok := kvCache[key]
		if !ok {
			return otto.UndefinedValue()
		}
		retV, _ := otto.ToValue(val)
		return retV
	})

	obj.Set("all", func(c otto.FunctionCall) otto.Value {
		kvObj, err := vm.Object("({})")
		for k, v := range kvCache {
			kvObj.Set(k, v)
		}
		return modules.ToResult(vm, kvObj, err)
	})

	obj.Set("load", func(c otto.FunctionCall) otto.Value {
		path, _ := c.Argument(0).ToString()
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return modules.ToResult(vm, nil, err)
		}
		err = json.Unmarshal(content, &kvCache)
		if err != nil {
			return modules.ToResult(vm, nil, err)
		}

		return otto.TrueValue()
	})

	obj.Set("save", func(c otto.FunctionCall) otto.Value {
		path, _ := c.Argument(0).ToString()
		content, _ := json.Marshal(kvCache)
		ioutil.WriteFile(path, content, 0777)
		return otto.TrueValue()
	})

	obj.Set("remove", func(c otto.FunctionCall) otto.Value {
		key, _ := c.Argument(0).ToString()
		delete(kvCache, key)
		return otto.TrueValue()
	})
	return obj.Value()
}

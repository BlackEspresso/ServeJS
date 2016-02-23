package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/robertkrimen/otto"
)

type FuncMapping map[string]func(w http.ResponseWriter, r *http.Request)
type JSCall func(otto.FunctionCall) otto.Value
type PluginInit func(*otto.Otto) otto.Value

type Plugin struct {
	Name        string
	Init        PluginInit
	Disabled    bool
	HttpMapping FuncMapping
}

type ModByName map[string]*Plugin

var modules ModByName = ModByName{}
var defaultPath string = "./js/main.js"
var usedRuntimes int = 0

func NewJSRuntime() (*otto.Otto, error) {
	vm := otto.New()
	RegisterModules(vm)
	path := defaultPath

	fileC, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	/*
		fileHash := md5.Sum(fileC)
		if len(lastMd5) == 0 || lastMd5 != fileHash {
			lastMd5 = fileHash
			fmt.Println("compiling")
			compiled, err = vm.Compile(path, nil)
			if err != nil {
				return nil, err
			}
		}
		_, err = vm.Run(compiled)
	*/
	_, err = vm.Run(string(fileC))

	if err == nil {
		usedRuntimes += 1
		runtime.SetFinalizer(vm, finalizer)
	}
	return vm, err
}

func finalizer(f *otto.Otto) {
	usedRuntimes -= 1
	fmt.Println("used runtimes ", usedRuntimes)
}

func RegisterModules(vm *otto.Otto) {
	vm.Set("require", func(c otto.FunctionCall) otto.Value {
		name, _ := c.Argument(0).ToString()
		module, ok := modules[name]
		if !ok {
			if strings.Index(name, "./") == 0 {
				c, err := ioutil.ReadFile(name)
				if err == nil {
					vm.Run(c)
				}
			} else {
				return otto.UndefinedValue()
			}
		}
		return module.Init(c.Otto)
	})
}

func AddPlugin(p *Plugin) {
	modules[p.Name] = p
}

func GetPlugins() []*Plugin {
	modulesList := []*Plugin{}
	for _, p := range modules {
		modulesList = append(modulesList, p)
	}
	return modulesList
}

func ToResult(vm *otto.Otto, valOk interface{}, err error) otto.Value {
	res, _ := vm.Object("({})")
	if err != nil {
		res.Set("error", err.Error())
	} else {
		res.Set("ok", valOk)
	}
	return res.Value()
}

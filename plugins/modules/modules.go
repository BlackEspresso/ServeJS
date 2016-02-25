package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/robertkrimen/otto"
)

type FuncMapping map[string]func(w http.ResponseWriter, r *http.Request)
type JSCall func(otto.FunctionCall) otto.Value
type PluginInit func(*JsVm) otto.Value

type Plugin struct {
	Name        string
	Init        PluginInit
	Disabled    bool
	HttpMapping FuncMapping
}

type ModByName map[string]*Plugin

var modules ModByName = ModByName{}
var defaultPath string = "./js/main.js"
var UsedRuntimes int = 0
var pluginName = "modules"

type JsVm struct {
	vm    *otto.Otto
	inUse *sync.Mutex
}

func NewJsVm(vm *otto.Otto) *JsVm {
	return &JsVm{vm, &sync.Mutex{}}
}

func (j *JsVm) Run(src interface{}) (otto.Value, error) {
	j.inUse.Lock()
	ret, err := j.vm.Run(src)
	j.inUse.Unlock()
	return ret, err
}

func (j *JsVm) Object(source string) (*otto.Object, error) {
	//j.inUse.Lock()
	ret, err := j.vm.Object(source)
	//j.inUse.Unlock()
	return ret, err
}

func (j *JsVm) ToValue(value interface{}) (otto.Value, error) {
	//j.inUse.Lock()
	val, err := j.vm.ToValue(value)
	//j.inUse.Unlock()
	return val, err
}

func (j *JsVm) Set(name string, value interface{}) error {
	//j.inUse.Lock()
	err := j.vm.Set(name, value)
	//j.inUse.Unlock()
	return err
}

func (j *JsVm) Call(source string, this interface{}, argumentList ...interface{}) (otto.Value, error) {
	j.inUse.Lock()
	ret, err := j.vm.Call(source, this, argumentList...)
	j.inUse.Unlock()
	return ret, err
}

func InitPlugin() *Plugin {

	p1 := Plugin{
		Name: pluginName,
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *JsVm) otto.Value {

	obj, _ := vm.Object("({})")

	obj.Set("run", func(c otto.FunctionCall) otto.Value {
		go func() {
			cFunc := c.ArgumentList[0]
			arg := make([]interface{}, 1)
			arg[0] = cFunc
			_, err := vm.Call(`Function.call.call`, nil, cFunc)
			fmt.Println(err)
		}()
		return otto.TrueValue()
	})

	return obj.Value()
}

func NewJSRuntime() (*JsVm, error) {
	vm := otto.New()
	jsvm := NewJsVm(vm)
	RegisterModules(jsvm)
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

	_, err = jsvm.Run(string(fileC))

	if err == nil {
		UsedRuntimes += 1
		runtime.SetFinalizer(jsvm, finalizer)
	}

	return jsvm, err
}

func finalizer(f *JsVm) {
	UsedRuntimes -= 1
	fmt.Println("used runtimes ", UsedRuntimes)
}

func RegisterModules(vm *JsVm) {
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
		return module.Init(vm)
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

func ToResult(vm *JsVm, valOk interface{}, err error) otto.Value {
	res, _ := vm.Object("({})")
	if err != nil {
		res.Set("error", err.Error())
	} else {
		res.Set("ok", valOk)
	}
	return res.Value()
}

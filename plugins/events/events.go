package events

import (
	"time"

	"./../modules"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "events",
		Init: registerVM,
	}

	return &p1
}

var channels map[string]chan interface{} = map[string]chan interface{}{}
var routes map[string][]string = map[string][]string{}

func registerVM(vm *modules.JsVm) otto.Value {
	k, _ := vm.Object("({})")

	k.Set("sleep", func(c otto.FunctionCall) otto.Value {
		ms, _ := c.Argument(0).ToInteger()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return otto.TrueValue()
	})

	k.Set("waitfor", func(c otto.FunctionCall) otto.Value {
		//vm.ToValue()
		return otto.TrueValue()
	})

	k.Set("route", func(c otto.FunctionCall) otto.Value {
		source, _ := c.Argument(0).ToString()
		target, _ := c.Argument(1).ToString()

		route, ok := routes[source]
		if !ok {
			route = []string{}
			route = append(route, target)
		}

		routes[source] = route
		return otto.TrueValue()
	})

	k.Set("next", func(c otto.FunctionCall) otto.Value {
		name, _ := c.Argument(0).ToString()
		stack, ok := channels[name]
		if !ok {
			return otto.UndefinedValue()
		}
		value := <-stack
		valueStr := ""
		if value == nil {
			return otto.UndefinedValue()
		} else {
			valueStr = value.(string)
		}
		vmObj, _ := otto.ToValue(valueStr)
		return vmObj
	})

	k.Set("create", func(c otto.FunctionCall) otto.Value {
		name, _ := c.Argument(0).ToString()
		GetChannel(name)
		return otto.TrueValue()
	})

	k.Set("push", func(c otto.FunctionCall) otto.Value {
		name, _ := c.Argument(0).ToString()
		value, _ := c.Argument(1).ToString()

		route, ok := routes[name]
		if ok {
			for _, k := range route {
				Push(k, value)
			}
		} else {
			Push(name, value)
		}

		return otto.TrueValue()
	})
	return k.Value()
}

func GetChannel(name string) chan interface{} {
	channel, ok := channels[name]
	if !ok {
		channel = make(chan interface{}, 5)
		channels[name] = channel
	}
	return channel
}

func Push(name, value string) {
	channel := GetChannel(name)
	channel <- value
}

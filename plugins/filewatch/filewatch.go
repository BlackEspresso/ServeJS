package filewatch

import (
	"log"

	"./../modules"
	"github.com/fsnotify/fsnotify"
	"github.com/robertkrimen/otto"
)

var pluginName string = "filewatch"

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: pluginName,
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *otto.Otto) otto.Value {
	watcher, err := fsnotify.NewWatcher()
	//defer watcher.Close()
	paths := map[string]bool{}
	if err != nil {
		log.Println(pluginName + " " + err.Error())
	}
	obj, _ := vm.Object("({})")

	obj.Set("watchDir", func(c otto.FunctionCall) otto.Value {
		path, _ := c.Argument(0).ToString()
		err := watcher.Add(path)
		paths[path] = true
		return modules.ToResult(vm, true, err)
	})
	obj.Set("start", func(c otto.FunctionCall) otto.Value {
		log.Println(pluginName + " start")
		go func() {
			for {
				select {
				case ev := <-watcher.Events:
					log.Println("event:", ev.String())
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()
		return otto.TrueValue()
	})
	return obj.Value()
}

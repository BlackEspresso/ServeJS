package stats

import (
	"runtime"

	"./../modules"
	"github.com/robertkrimen/otto"
)

var pluginName string = "stats"

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: pluginName,
		Init: registerVM,
	}

	return &p1
}

type StatsObj struct {
	UsedVMs          int
	MemoryAlloc      uint64
	MemoryTotalAlloc uint64
}

func registerVM(vm *modules.JsVm) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("getStats", func(c otto.FunctionCall) otto.Value {
		memStat := runtime.MemStats{}
		runtime.ReadMemStats(&memStat)

		var stat = StatsObj{
			modules.UsedRuntimes,
			memStat.Alloc,
			memStat.TotalAlloc,
		}

		val, _ := vm.ToValue(stat)
		return val
	})

	return obj.Value()
}

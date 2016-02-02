package tasks

import (
	"time"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

type TaskFunc func(*TaskBlock)

type TaskBlock struct {
	Name        string
	Start       int64
	Repeat      int
	RunTime     int
	Done        bool
	Success     bool
	SuccessText string
	ErrorText   string
	Func        TaskFunc `json:"-"`
}

var tasks []*TaskBlock = []*TaskBlock{}

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "tasks",
		Init: func(vm *otto.Otto) {
			vm.Set("addTask", func(c otto.FunctionCall) otto.Value {
				t := TaskBlock{}
				t.Name, _ = c.Argument(0).ToString()
				t.Start, _ = c.Argument(1).ToInteger()
				fs := c.Argument(2)
				t.Func = func(tb *TaskBlock) {
					fs.Call(otto.Value{})
				}
				tasks = append(tasks, &t)
				return otto.TrueValue()
			})
			vm.Set("startTasks", func(c otto.FunctionCall) otto.Value {
				Start()
				return otto.TrueValue()
			})
		},
	}

	return &p1
}

func RunTasks() {
	for _, v := range tasks {
		if !v.Done {
			v.Func(v)
		}
	}
}

func Start() {
	go func() {
		for {
			RunTasks()
			time.Sleep(time.Second * 5)
		}
	}()
}

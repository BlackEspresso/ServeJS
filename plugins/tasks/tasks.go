package tasks

import (
	"log"
	"time"

	"./../modules"
	"github.com/robertkrimen/otto"
)

type TaskFunc func(*TaskBlock)

type TaskBlock struct {
	Id          int
	Name        string
	Start       time.Time
	Repeat      int
	RunTime     int
	Done        bool
	Success     bool
	SuccessText string
	ErrorText   string
	Func        TaskFunc `json:"-"`
}

var tasks []*TaskBlock = []*TaskBlock{}

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "tasks",
		Init: func(vm *modules.JsVm) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("addTask", func(c otto.FunctionCall) otto.Value {
				t := TaskBlock{}
				t.Name, _ = c.Argument(0).ToString()
				startTimeStr := c.Argument(1).String()
				startTime, _ := time.Parse(time.RFC1123, startTimeStr)
				t.Start = startTime
				fs := c.Argument(2)
				t.Func = func(tb *TaskBlock) {
					val, err := fs.Call(otto.Value{})
					tb.Done = true
					if err != nil {
						tb.ErrorText = err.Error()
						log.Println(err)
					}
					tb.SuccessText, _ = val.ToString()
					if tb.Repeat > 0 {
						tb.Start.Add(time.Second * time.Duration(tb.Repeat))
					}
				}
				tasks = append(tasks, &t)
				return otto.TrueValue()
			})

			obj.Set("startTasks", func(c otto.FunctionCall) otto.Value {
				Start()
				return otto.TrueValue()
			})
			return obj.Value()
		},
	}

	return &p1
}

func RunTasks() {
	for _, v := range tasks {
		if !v.Done && time.Now().After(v.Start) {
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

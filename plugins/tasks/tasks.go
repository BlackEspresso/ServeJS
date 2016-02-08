package tasks

import (
	"fmt"
	"time"

	"./../pluginbase"
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

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "tasks",
		Init: func(vm *otto.Otto) {
			vm.Set("addTask", func(c otto.FunctionCall) otto.Value {
				t := TaskBlock{}
				t.Name, _ = c.Argument(0).ToString()
				startTimeStr := c.Argument(1).String()
				startTime, _ := time.Parse(time.RFC1123, startTimeStr)
				fmt.Println(startTimeStr, startTime)
				t.Start = startTime
				fs := c.Argument(2)
				t.Func = func(tb *TaskBlock) {
					val, err := fs.Call(otto.Value{})
					tb.Done = true
					if err != nil {
						tb.ErrorText = err.Error()
					}
					tb.SuccessText, _ = val.ToString()
					if tb.Repeat > 0 {
						tb.Start.Add(time.Second * time.Duration(tb.Repeat))
					}
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

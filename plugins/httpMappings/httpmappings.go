package httpmappings

import (
	"net/http"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

type ServeFunc func(w http.ResponseWriter, r *http.Request)

var urlMapping map[string]string = map[string]string{}

func InitPlugin() *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "httpmapping",
		Init: func(vm *otto.Otto) {
			vm.Set("addMapping", func(c otto.FunctionCall) otto.Value {
				url, _ := c.Argument(0).ToString()
				name, _ := c.Argument(1).ToString()
				AddMapping(url, name)
				return otto.TrueValue()
			})
		},
	}

	return &p
}

func AddMapping(url string, name string) {
	urlMapping[url] = name
}

func RunMappings(w http.ResponseWriter, r *http.Request, plugins []*pluginbase.Plugin) bool {
	url := r.URL.Path
	funcName, ok := urlMapping[url]
	if ok {
		for _, p := range plugins {
			f, hasFunc := p.HttpMapping[funcName]
			if hasFunc {
				f(w, r)
				return true
			}
		}
	}
	return false
}

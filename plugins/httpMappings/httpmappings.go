package httpmappings

import (
	"io/ioutil"
	"log"
	"net/http"

	"./../modules"
	"github.com/robertkrimen/otto"
)

type ServeFunc func(w http.ResponseWriter, r *http.Request)

var urlMapping map[string]string = map[string]string{}

func InitPlugin() *modules.Plugin {
	p := modules.Plugin{
		Name: "httpmappings",
		Init: func(vm *otto.Otto) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("addMapping", func(c otto.FunctionCall) otto.Value {
				url, _ := c.Argument(0).ToString()
				name, _ := c.Argument(1).ToString()
				AddMapping(url, name)
				return otto.TrueValue()
			})
			return obj.Value()
		},
		HttpMapping: modules.FuncMapping{
			"writefile": func(w http.ResponseWriter, r *http.Request) {
				writeMainJs(w, r)
			},
		},
	}

	return &p
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func writeMainJs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mainjs := r.FormValue("mainjs")
	if mainjs != "" {
		fileContent, err := ioutil.ReadFile("./js/main.js")
		logError(err)
		err = ioutil.WriteFile("./js/main.js.bak", fileContent, 0777)
		logError(err)
	}
	ioutil.WriteFile("./js/main.js", []byte(mainjs), 0777)
}

func AddMapping(url string, name string) {
	urlMapping[url] = name
}

func RunMappings(w http.ResponseWriter, r *http.Request, plugins []*modules.Plugin) bool {
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

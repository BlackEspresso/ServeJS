package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"./plugins/cache"
	"./plugins/cmd"
	"./plugins/dns"
	"./plugins/file"
	"./plugins/htmlcheck"
	phttp "./plugins/http"
	"./plugins/httpmappings"
	"./plugins/mail"
	"./plugins/modules"
	"./plugins/tasks"
	"./plugins/templating"
	"./plugins/websocket"

	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port    int
	Plugins map[string]map[string]string
}

var serverConf = Configuration{}

func main() {
	configraw, err := ioutil.ReadFile("./serverjs.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configraw, &serverConf)
	if err != nil {
		log.Fatal(err)
	}

	registerPlugins()

	vm, err := newJSRuntime()
	if err != nil {
		log.Fatal(err)
	}

	_, err = vm.Call("onServerStart", nil)
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/", jsHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func registerPlugins() {
	modules.AddPlugin(templating.InitPlugin())
	modules.AddPlugin(file.InitPlugin())
	modules.AddPlugin(dns.InitPlugin())
	modules.AddPlugin(tasks.InitPlugin())
	modules.AddPlugin(mail.InitPlugin())
	modules.AddPlugin(cmd.InitPlugin())
	modules.AddPlugin(phttp.InitPlugin())
	modules.AddPlugin(websocket.InitPlugin(newJSRuntime))
	modules.AddPlugin(httpmappings.InitPlugin())
	modules.AddPlugin(cache.InitPlugin())
	modules.AddPlugin(htmlcheck.InitPlugin())
}

func newJSRuntime() (*otto.Otto, error) {
	vm := otto.New()
	modules.RegisterRequire(vm)

	vm.Set("settings", serverConf)
	fileC, err := ioutil.ReadFile("./js/main.js")
	if err != nil {
		return nil, err
	}
	_, err = vm.Run(string(fileC))

	return vm, err
}

func jsHandler(w http.ResponseWriter, r *http.Request) {

	ret := httpmappings.RunMappings(w, r, modules.GetPlugins())
	if ret {
		// httpmappings has process this mapping
		return
	}

	vm, err := newJSRuntime()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	objRequest, _ := vm.Object("({})")
	objResponse, _ := vm.Object("({})")
	phttp.RequestToJso(objRequest, r)
	phttp.ResponseWriterToJso(objResponse, w)

	_, err = vm.Call("onRequest", nil, objResponse, objRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	phttp.JsoToResponseWriter(objResponse, w)
}

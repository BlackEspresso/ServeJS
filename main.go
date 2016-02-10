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
	"./plugins/pluginbase"
	"./plugins/tasks"
	"./plugins/templating"
	"./plugins/websocket"

	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

var plugins []*pluginbase.Plugin = []*pluginbase.Plugin{}

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

	addPlugins()

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

func addPlugins() {
	p := templating.InitPlugin()
	addPlugin(p)
	p = file.InitPlugin()
	addPlugin(p)
	p = dns.InitPlugin()
	addPlugin(p)
	p = tasks.InitPlugin()
	addPlugin(p)
	p = mail.InitPlugin()
	addPlugin(p)
	p = cmd.InitPlugin()
	addPlugin(p)
	p = phttp.InitPlugin()
	addPlugin(p)
	p = phttp.InitPlugin()
	addPlugin(p)
	p = websocket.InitPlugin(newJSRuntime)
	addPlugin(p)
	p = httpmappings.InitPlugin()
	addPlugin(p)
	p = cache.InitPlugin()
	addPlugin(p)
	p = htmlcheck.InitPlugin()
	addPlugin(p)
}

func newJSRuntime() (*otto.Otto, error) {
	vm := otto.New()
	for _, v := range plugins {
		v.Init(vm)
	}

	vm.Set("settings", serverConf)
	fileC, err := ioutil.ReadFile("./js/main.js")
	if err != nil {
		return nil, err
	}
	_, err = vm.Run(string(fileC))

	return vm, err
}

func jsHandler(w http.ResponseWriter, r *http.Request) {

	ret := httpmappings.RunMappings(w, r, plugins)
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

func addPlugin(p *pluginbase.Plugin) {
	plugins = append(plugins, p)
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"./plugins/cmd"
	"./plugins/dns"
	"./plugins/file"
	phttp "./plugins/http"
	"./plugins/mail"
	"./plugins/pluginbase"
	"./plugins/tasks"
	"./plugins/templating"
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
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path[1:])
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/", jsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	vm := newJSRuntime(r)
	fileC, err := ioutil.ReadFile("./js/main.js")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	objRequest, _ := vm.Object("({})")
	objResponse, _ := vm.Object("({})")
	phttp.RequestToJSObject(objRequest, r)
	phttp.ResponseWriterToJSObject(objResponse, w)

	_, err = vm.Run(string(fileC))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = vm.Call("onRequest", nil, objResponse, objRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	phttp.JsObjectToResponse(objResponse, w)

	//str, _ := ret.ToString()
	//w.Write([]byte(str))
}

func newJSRuntime(r *http.Request) *otto.Otto {
	vm := otto.New()
	for _, v := range plugins {
		v.Init(vm)
	}

	vm.Set("settings", serverConf)
	return vm
}

func addPlugin(p *pluginbase.Plugin) {
	plugins = append(plugins, p)
}

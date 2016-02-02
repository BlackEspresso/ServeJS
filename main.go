package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./plugins/dns"
	"./plugins/file"
	"./plugins/pluginbase"
	"./plugins/tasks"
	"./plugins/templating"
	"github.com/robertkrimen/otto"
)

var plugins []*pluginbase.Plugin = []*pluginbase.Plugin{}

func main() {
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
	requestToJSObject(objRequest, r)
	responseToJSObject(objResponse, w)

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

	jsObjectToResponse(objResponse, w)

	//str, _ := ret.ToString()
	//w.Write([]byte(str))
}

func newJSRuntime(r *http.Request) *otto.Otto {
	vm := otto.New()
	for _, v := range plugins {
		v.Init(vm)
	}

	return vm
}

func addPlugin(p *pluginbase.Plugin) {
	plugins = append(plugins, p)
}

func responseToJSObject(o *otto.Object, w http.ResponseWriter) {
	o.Set("write", func(c otto.FunctionCall) otto.Value {
		text, _ := c.Argument(0).ToString()
		w.Write([]byte(text))
		return otto.TrueValue()
	})
}

func jsObjectToResponse(respObj *otto.Object, w http.ResponseWriter) {
	contentTypeV, err := respObj.Get("contentType")
	if err == nil {
		contentType, _ := contentTypeV.ToString()
		w.Header().Set("Content-Type", contentType)
	}

	codeV, err := respObj.Get("statusCode")
	if err == nil && codeV.IsDefined() {
		code, _ := codeV.ToInteger()
		fmt.Println(code)
		w.WriteHeader(int(code))
	}
}

func requestToJSObject(o *otto.Object, r *http.Request) {
	o.Set("url", r.URL.String())
	o.Set("header", r.Header)
	o.Set("cookies", r.Cookies())
	o.Set("method", r.Method)
	r.ParseForm()
	o.Set("formValues", r.Form)
}

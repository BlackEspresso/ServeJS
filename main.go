package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robertkrimen/otto"
)

var templates *template.Template

func main() {
	templates, _ = template.ParseGlob("./tmpl/*.thtml")
	http.HandleFunc("/", jsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	objResponse.Set("write", func(c otto.FunctionCall) otto.Value {
		text, _ := c.Argument(0).ToString()
		w.Write([]byte(text))
		return otto.TrueValue()
	})

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

	contentTypeV, err := objResponse.Get("contentType")
	if err == nil {
		contentType, _ := contentTypeV.ToString()
		w.Header().Set("Content-Type", contentType)
	}

	codeV, err := objResponse.Get("statusCode")
	if err == nil && codeV.IsDefined() {
		code, _ := codeV.ToInteger()
		fmt.Println(code)
		w.WriteHeader(int(code))
	}

	//str, _ := ret.ToString()
	//w.Write([]byte(str))
}

type M struct {
	Name string
}

func newJSRuntime(r *http.Request) *otto.Otto {
	vm := otto.New()
	vm.Set("settings", 3)
	vm.Set("reloadTemplates", func(c otto.FunctionCall) otto.Value {
		templates, _ = template.ParseGlob("./tmpl/*.thtml")
		return otto.TrueValue()
	})
	vm.Set("template", func(c otto.FunctionCall) otto.Value {
		name, _ := c.Argument(0).ToString()
		text, _ := c.Argument(1).ToString()

		b := new(bytes.Buffer)
		t := templates.Lookup(name)
		if t == nil {
			return otto.UndefinedValue()
		}
		t.Parse(name)

		err := t.Execute(b, text)
		if err != nil {
			return otto.UndefinedValue()
		}
		retV, _ := otto.ToValue(b.String())
		return retV
	})
	return vm
}

func requestToJSObject(o *otto.Object, r *http.Request) *otto.Object {
	o.Set("url", r.URL.String())
	o.Set("header", r.Header)
	o.Set("cookies", r.Cookies())
	o.Set("method", r.Method)
	return o
}

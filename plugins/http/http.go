package http

import (
	"io/ioutil"
	"net/http"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "file",
		Init: func(vm *otto.Otto) {
			vm.Set("httpDo", func(c otto.FunctionCall) otto.Value {
				arg1 := c.Argument(0)
				req := JSObjectToRequest(arg1.Object())

				client := &http.Client{}
				resp, _ := client.Do(req)

				respObj, _ := vm.Object("({})")
				ResponseToJSObject(respObj, resp)
				return respObj.Value()
			})

		},
	}
	return &p
}

func ResponseToJSObject(o *otto.Object, w *http.Response) {
	o.Set("status", w.Status)
	o.Set("header", w.Header)

	c, _ := ioutil.ReadAll(w.Body)
	o.Set("body", string(c))
}

func ResponseWriterToJSObject(o *otto.Object, w http.ResponseWriter) {
	o.Set("write", func(c otto.FunctionCall) otto.Value {
		text, _ := c.Argument(0).ToString()
		w.Write([]byte(text))
		return otto.TrueValue()
	})
}

func JsObjectToResponse(respObj *otto.Object, w http.ResponseWriter) {
	contentTypeV, err := respObj.Get("contentType")
	if err == nil {
		contentType, _ := contentTypeV.ToString()
		w.Header().Set("Content-Type", contentType)
	}

	codeV, err := respObj.Get("statusCode")
	if err == nil && codeV.IsDefined() {
		code, _ := codeV.ToInteger()
		w.WriteHeader(int(code))
	}
}

func JSObjectToRequest(o *otto.Object) *http.Request {
	//var buffer := bytes.NewReader()
	url, _ := o.Get("url")
	//header, _ := o.Get("header")
	method, _ := o.Get("method")

	urlStr, _ := url.ToString()
	//headerStr, _ := header.ToString()
	methodStr, _ := method.ToString()

	req, _ := http.NewRequest(methodStr, urlStr, nil)
	return req
}

func RequestToJSObject(o *otto.Object, r *http.Request) {
	o.Set("url", r.URL.String())
	o.Set("header", r.Header)
	o.Set("cookies", r.Cookies())
	o.Set("method", r.Method)
	r.ParseForm()
	o.Set("formValues", r.Form)
}

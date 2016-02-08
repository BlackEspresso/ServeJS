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
		Init: registerVM,
	}
	return &p
}

func registerVM(vm *otto.Otto) {
	obj, _ := vm.Object("({})")
	vm.Set("http", obj)

	obj.Set("do", func(c otto.FunctionCall) otto.Value {
		arg1 := c.Argument(0)
		req := JsoToRequest(arg1.Object())

		client := &http.Client{}
		resp, _ := client.Do(req)

		respObj, _ := vm.Object("({})")
		ResponseToJso(respObj, resp)
		return respObj.Value()
	})
}

func ResponseToJso(o *otto.Object, w *http.Response) {
	o.Set("status", w.Status)
	o.Set("header", w.Header)
	o.Set("cookies", w.Cookies())
	o.Set("statusCode", w.StatusCode)
	o.Set("proto", w.Proto)
	c, _ := ioutil.ReadAll(w.Body)
	o.Set("body", string(c))
}

func RequestToJso(o *otto.Object, r *http.Request) {
	o.Set("url", r.URL.String())
	o.Set("header", r.Header)
	o.Set("cookies", r.Cookies())
	o.Set("method", r.Method)
	r.ParseForm()
	o.Set("formValues", r.Form)
}

func ResponseWriterToJso(o *otto.Object, w http.ResponseWriter) {
	o.Set("write", func(c otto.FunctionCall) otto.Value {
		text, _ := c.Argument(0).ToString()
		w.Write([]byte(text))
		return otto.TrueValue()
	})
}

func JsoToResponseWriter(respObj *otto.Object, w http.ResponseWriter) {
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

func JsoToRequest(o *otto.Object) *http.Request {
	//var buffer := bytes.NewReader()
	url, _ := o.Get("url")
	header, _ := o.Get("header")
	method, err := o.Get("method")
	methodStr, err := method.ToString()
	if method == otto.UndefinedValue() || err != nil {
		methodStr = "GET"
	}

	urlStr, _ := url.ToString()
	headerIface, err := header.Export()

	req, _ := http.NewRequest(methodStr, urlStr, nil)
	headerMap, ok := headerIface.(map[string]string)
	if ok {
		for k, h := range headerMap {
			req.Header.Add(k, h)
		}
	}
	return req
}

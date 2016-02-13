package httplistener

import (
	"log"
	"net/http"

	phttp "./../http"
	"./../httpmappings"
	"./../modules"
	"github.com/robertkrimen/otto"
)

var kvCache map[string]string = map[string]string{}

func InitPlugin(createVM func() (*otto.Otto, error)) *modules.Plugin {

	p1 := modules.Plugin{
		Name: "httplistener",
		Init: func(vm *otto.Otto) otto.Value {
			obj, _ := vm.Object("({})")
			obj.Set("addr", ":8081")
			obj.Set("start", func(c otto.FunctionCall) otto.Value {
				addrObj, _ := obj.Get("addr")
				addr, _ := addrObj.ToString()

				go startServer(addr, createVM)
				return otto.TrueValue()
			})
			obj.Set("startAndWait", func(c otto.FunctionCall) otto.Value {
				addrObj, _ := obj.Get("addr")
				addr, _ := addrObj.ToString()

				startServer(addr, createVM)
				return otto.TrueValue()
			})
			return obj.Value()
		},
	}

	return &p1
}

func startServer(addr string, newJSRuntime func() (*otto.Otto, error)) {
	jss := &JsServer{}
	jss.vmFunc = newJSRuntime

	log.Fatal(http.ListenAndServe(addr, jss))
}

type JsServer struct {
	vmFunc func() (*otto.Otto, error)
}

func (j *JsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ret := httpmappings.RunMappings(w, r, modules.GetPlugins())
	if ret {
		// httpmappings has process this mapping
		return
	}

	vm, err := j.vmFunc()
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
		return
	}

	phttp.JsoToResponseWriter(objResponse, w)

}

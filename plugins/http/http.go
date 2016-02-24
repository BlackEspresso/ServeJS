package http

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"./../modules"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {
	p := modules.Plugin{
		Name: "http",
		Init: registerVM,
		HttpMapping: modules.FuncMapping{
			"servefile": func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.URL.Path)
				http.ServeFile(w, r, r.URL.Path[1:])
			},
		},
	}
	return &p
}

func registerVM(vm *modules.JsVm) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("do", func(c otto.FunctionCall) otto.Value {
		arg1 := c.Argument(0)
		req := JsoToRequest(arg1.Object())
		followRedirObj, err := arg1.Object().Get("followRedirects")
		followRedirects := true
		if followRedirObj != otto.UndefinedValue() {
			followRedirects, _ = followRedirObj.ToBoolean()
		}

		client := &http.Client{}

		var resp *http.Response
		if followRedirects {
			resp, err = client.Do(req)
		} else {
			resp, err = http.DefaultTransport.RoundTrip(req)
		}

		respObj, _ := vm.Object("({})")
		if err == nil {
			ResponseToJso(respObj, resp)
		}

		return modules.ToResult(vm, respObj, err)
	})
	return obj.Value()
}

func ResponseToJso(o *otto.Object, w *http.Response) {
	o.Set("status", w.Status)
	o.Set("header", w.Header)
	o.Set("cookies", w.Cookies())
	o.Set("statusCode", w.StatusCode)
	o.Set("proto", w.Proto)
	c, _ := ioutil.ReadAll(w.Body)
	o.Set("body", string(c))
	if w.TLS != nil {
		o.Set("tlsServerName", w.TLS.ServerName)
		o.Set("tlsNegotiatedProtocol", w.TLS.NegotiatedProtocol)
		names := []string{}
		ips := []string{}
		for _, k := range w.TLS.PeerCertificates {
			for _, n := range k.DNSNames {
				names = append(names, n)
			}
			for _, n := range k.IPAddresses {
				ips = append(ips, n.String())
			}
		}
		o.Set("tlsDNSNames", names)
		o.Set("tlsIPs", ips)
	}
}

func RequestToJso(o *otto.Object, r *http.Request) {
	o.Set("url", r.URL.String())
	o.Set("header", r.Header)
	o.Set("cookies", r.Cookies())
	o.Set("method", r.Method)
	o.Set("host", r.Host)
	o.Set("contentLength", r.ContentLength)
	o.Set("proto", r.Proto)
	o.Set("transferEncoding", r.TransferEncoding)
	r.ParseForm()
	o.Set("formValues", r.Form)
}

func ResponseWriterToJso(o *otto.Object, w http.ResponseWriter) {
	o.Set("write", func(c otto.FunctionCall) otto.Value {
		text, _ := c.Argument(0).ToString()
		w.Write([]byte(text))
		return otto.TrueValue()
	})
	o.Set("writeHeader", func(c otto.FunctionCall) otto.Value {
		statusCode, _ := c.Argument(0).ToInteger()
		w.WriteHeader(int(statusCode))
		return otto.TrueValue()
	})
}

func JsoToResponseWriter(respObj *otto.Object, w http.ResponseWriter) {
	contentTypeV, err := respObj.Get("contentType")
	if err == nil {
		contentType, _ := contentTypeV.ToString()
		w.Header().Set("Content-Type", contentType)
	}

	wHeader := w.Header()
	setHeader(respObj, &wHeader)

	codeV, err := respObj.Get("statusCode")
	if err == nil && codeV.IsDefined() {
		code, _ := codeV.ToInteger()
		w.WriteHeader(int(code))
	}
}

func setHeader(o *otto.Object, h *http.Header) {
	header, _ := o.Get("header")
	headerIface, _ := header.Export()
	headerMap, ok := headerIface.(map[string]interface{})

	if ok {
		for k, v := range headerMap {
			strV, ok := v.(string)
			if ok {
				h.Set(k, strV)
			}
		}
	}
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func JsoToRequest(o *otto.Object) *http.Request {
	url, _ := o.Get("url")
	method, err := o.Get("method")
	body, _ := o.Get("body")

	methodStr, err := method.ToString()
	if method == otto.UndefinedValue() || err != nil {
		methodStr = "GET"
	}
	urlStr, _ := url.ToString()

	req, _ := http.NewRequest(methodStr, urlStr, nil)

	if body != otto.UndefinedValue() {
		str, err := body.ToString()
		if err == nil {
			req.Body = &ClosingBuffer{bytes.NewBufferString(str)}
		}
	}

	setHeader(o, &req.Header)

	return req
}

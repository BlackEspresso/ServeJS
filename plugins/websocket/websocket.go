package websocket

import (
	"fmt"
	"log"
	"net/http"

	"./../modules"
	"github.com/gorilla/websocket"
	"github.com/robertkrimen/otto"
)

var upgrader = websocket.Upgrader{}

func InitPlugin(createVM func() (*modules.JsVm, error)) *modules.Plugin {
	p := modules.Plugin{
		Name: "websocket",
		Init: func(vm *modules.JsVm) otto.Value {
			return otto.UndefinedValue()
		},
		HttpMapping: modules.FuncMapping{
			"websocket": func(w http.ResponseWriter, r *http.Request) {
				doWebSocket(w, r, createVM)
			},
		},
	}

	return &p
}

func doWebSocket(w http.ResponseWriter, r *http.Request,
	createVM func() (*modules.JsVm, error)) {

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	vm, err := createVM()

	fmt.Println("websocket start")

	obj, _ := vm.Object("({})")
	obj.Set("read", func(c otto.FunctionCall) otto.Value {
		_, message, _ := connection.ReadMessage()
		val, _ := otto.ToValue(message)
		return val
	})

	obj.Set("write", func(c otto.FunctionCall) otto.Value {
		mType, _ := c.Argument(0).ToInteger()
		message, _ := c.Argument(1).ToString()

		err := connection.WriteMessage(int(mType), []byte(message))
		val, _ := otto.ToValue(err)
		return val
	})

	_, err = vm.Call("onWebSocket", nil, obj.Value())
	if err != nil {
		log.Println("jserror", err)
	}

	fmt.Println("websocket end")
}

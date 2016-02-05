package websocket

import (
	"log"
	"net/http"

	"./../pluginbase"
	"github.com/gorilla/websocket"
	"github.com/robertkrimen/otto"
)

var upgrader = websocket.Upgrader{}

func InitPlugin(createVM func() *otto.Otto) *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "file",
		Init: func(vm *otto.Otto) {

			vm.Set("newWebsocket", func(c otto.FunctionCall) otto.Value {
				return otto.TrueValue()
			})

		},
	}
	return &p
}

func doWebSocket(w http.ResponseWriter, r *http.Request, createVM func() *otto.Otto) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	//vm := createVM()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

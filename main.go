package main

import (
	"io/ioutil"
	"log"

	"./plugins/cache"
	"./plugins/cmd"
	"./plugins/dns"
	"./plugins/file"
	"./plugins/goquery"
	"./plugins/htmlcheck"
	phttp "./plugins/http"
	"./plugins/httplistener"
	"./plugins/httpmappings"
	"./plugins/mail"
	"./plugins/modules"
	"./plugins/mongodb"
	"./plugins/settings"
	"./plugins/tasks"
	"./plugins/templating"
	"./plugins/time"
	"./plugins/websocket"

	"github.com/robertkrimen/otto"
)

func main() {
	registerPlugins()

	vm, err := newJSRuntime()
	if err != nil {
		log.Fatal(err)
	}

	_, err = vm.Call("onStart", nil)
	if err != nil {
		log.Println(err)
	}
}

func registerPlugins() {
	modules.AddPlugin(templating.InitPlugin())
	modules.AddPlugin(file.InitPlugin())
	modules.AddPlugin(dns.InitPlugin())
	modules.AddPlugin(tasks.InitPlugin())
	modules.AddPlugin(mail.InitPlugin())
	modules.AddPlugin(cmd.InitPlugin())
	modules.AddPlugin(phttp.InitPlugin())
	modules.AddPlugin(websocket.InitPlugin(newJSRuntime))
	modules.AddPlugin(httpmappings.InitPlugin())
	modules.AddPlugin(httplistener.InitPlugin(newJSRuntime))
	modules.AddPlugin(cache.InitPlugin())
	modules.AddPlugin(htmlcheck.InitPlugin())
	modules.AddPlugin(settings.InitPlugin())
	modules.AddPlugin(goquery.InitPlugin())
	modules.AddPlugin(mongodb.InitPlugin())
	modules.AddPlugin(time.InitPlugin())
}

func newJSRuntime() (*otto.Otto, error) {
	vm := otto.New()
	modules.RegisterRequire(vm)

	fileC, err := ioutil.ReadFile("./js/main.js")
	if err != nil {
		return nil, err
	}
	_, err = vm.Run(string(fileC))

	return vm, err
}

package main

import (
	//"crypto/md5"
	//"fmt"
	"log"

	"./plugins/cache"
	"./plugins/cmd"
	"./plugins/crypto"
	"./plugins/dns"
	"./plugins/file"
	"./plugins/filewatch"
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

	vm, err := modules.NewJSRuntime()
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
	modules.AddPlugin(websocket.InitPlugin(modules.NewJSRuntime))
	modules.AddPlugin(httpmappings.InitPlugin())
	modules.AddPlugin(httplistener.InitPlugin(modules.NewJSRuntime))
	modules.AddPlugin(cache.InitPlugin())
	modules.AddPlugin(htmlcheck.InitPlugin())
	modules.AddPlugin(settings.InitPlugin())
	modules.AddPlugin(goquery.InitPlugin())
	modules.AddPlugin(mongodb.InitPlugin())
	modules.AddPlugin(time.InitPlugin())
	modules.AddPlugin(crypto.InitPlugin())
	modules.AddPlugin(filewatch.InitPlugin())
}

var lastMd5 [16]byte = [16]byte{}
var compiled *otto.Script

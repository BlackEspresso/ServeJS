package main

import (
	"io/ioutil"
	"log"

	"./plugins/cache"
	"./plugins/cmd"
	"./plugins/dns"
	"./plugins/file"
	"./plugins/htmlcheck"
	phttp "./plugins/http"
	"./plugins/httplistener"
	"./plugins/httpmappings"
	"./plugins/mail"
	"./plugins/modules"
	"./plugins/tasks"
	"./plugins/templating"
	"./plugins/websocket"

	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port    int
	Plugins map[string]map[string]string
}

var serverConf = Configuration{}

func main() {
	configraw, err := ioutil.ReadFile("./serverjs.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configraw, &serverConf)
	if err != nil {
		log.Fatal(err)
	}

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
}

func newJSRuntime() (*otto.Otto, error) {
	vm := otto.New()
	modules.RegisterRequire(vm)

	vm.Set("settings", serverConf)
	fileC, err := ioutil.ReadFile("./js/main.js")
	if err != nil {
		return nil, err
	}
	_, err = vm.Run(string(fileC))

	return vm, err
}

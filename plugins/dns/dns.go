package dns

import (
	"net"

	"./../pluginbase"
	"github.com/miekg/dns"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "dns",
		Init: func(vm *otto.Otto) {
			vm.Set("resolve", func(c otto.FunctionCall) otto.Value {
				name, err := c.Argument(0).ToString()
				if err != nil {
					return otto.UndefinedValue()
				}

				config, _ := dns.ClientConfigFromFile("./resolv.conf")
				client := new(dns.Client)

				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(name), dns.TypeANY)
				m.RecursionDesired = true

				host := net.JoinHostPort(config.Servers[0], config.Port)
				r, _, err := client.Exchange(m, host)
				if err != nil {
					k, _ := otto.ToValue("")
					return k
				}
				resp := ""
				for _, v := range r.Answer {
					resp += v.String()
					resp += "\n"
				}
				respV, _ := otto.ToValue(resp)
				return respV
			})
		},
	}

	return &p1
}

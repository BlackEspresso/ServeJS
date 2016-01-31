package templating

import (
	"bytes"
	"html/template"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

var templates *template.Template

func InitPlugin() *pluginbase.Plugin {
	templates, _ = template.ParseGlob("./tmpl/*.thtml")

	p1 := pluginbase.Plugin{
		Name: "template",
		Init: func(vm *otto.Otto) {
			vm.Set("template", loadTemplate)
			vm.Set("reloadTemplate", reloadTemplate)
		},
	}

	return &p1
}

func reloadTemplate(c otto.FunctionCall) otto.Value {
	templates, _ = template.ParseGlob("./tmpl/*.thtml")
	return otto.TrueValue()
}

func loadTemplate(c otto.FunctionCall) otto.Value {
	name, _ := c.Argument(0).ToString()
	text, _ := c.Argument(1).ToString()

	b := new(bytes.Buffer)
	t := templates.Lookup(name)
	if t == nil {
		return otto.UndefinedValue()
	}
	t.Parse(name)

	err := t.Execute(b, text)
	if err != nil {
		return otto.UndefinedValue()
	}
	retV, _ := otto.ToValue(b.String())
	return retV
}

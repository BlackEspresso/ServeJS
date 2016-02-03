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
			vm.Set("runTemplate", func(c otto.FunctionCall) otto.Value {
				name, _ := c.Argument(0).ToString()
				jsObject := c.Argument(1).Object()

				b := new(bytes.Buffer)
				t := templates.Lookup(name)
				if t == nil {
					return otto.UndefinedValue()
				}
				t.Parse(name)

				kv := map[string]interface{}{}
				for _, k := range jsObject.Keys() {
					objValue, _ := jsObject.Get(k)
					kv[k] = objValue
				}

				err := t.Execute(b, kv)
				if err != nil {
					return otto.UndefinedValue()
				}
				retV, _ := otto.ToValue(b.String())
				return retV
			})

			vm.Set("reloadTemplates", reloadTemplates)
		},
	}

	return &p1
}

func reloadTemplates(c otto.FunctionCall) otto.Value {
	templates, _ = template.ParseGlob("./tmpl/*.thtml")
	return otto.TrueValue()
}

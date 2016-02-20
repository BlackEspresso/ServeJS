package htmlcheck

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BlackEspresso/htmlcheck"

	"./../modules"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "htmlcheck",
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *otto.Otto) otto.Value {
	validater := htmlcheck.Validator{}
	obj, _ := vm.Object("({})")

	obj.Set("loadTags", func(c otto.FunctionCall) otto.Value {
		path, err := c.Argument(0).ToString()
		tags, err := LoadTagsFromFile(path)
		if err == nil {
			validater.AddValidTags(tags)
		}
		return modules.ToResult(vm, true, err)
	})
	obj.Set("validate", func(c otto.FunctionCall) otto.Value {
		htmltext, _ := c.Argument(0).ToString()
		errors := validater.ValidateHtmlString(htmltext)
		objs, _ := vm.ToValue(errors)
		return objs
	})
	return obj.Value()
}

func LoadTagsFromFile(path string) ([]*htmlcheck.ValidTag, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []*htmlcheck.ValidTag{}, err
	}

	var validTags []*htmlcheck.ValidTag
	err = json.Unmarshal(content, &validTags)
	if err != nil {
		return []*htmlcheck.ValidTag{}, err
	}

	return validTags, nil
}

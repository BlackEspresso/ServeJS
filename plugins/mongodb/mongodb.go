package mongodb

import (
	"./../modules"
	"github.com/robertkrimen/otto"
	"gopkg.in/mgo.v2"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "mongodb",
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *modules.JsVm) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("newSession", func(c otto.FunctionCall) otto.Value {
		str, _ := c.Argument(0).ToString()
		session, err := mgo.Dial(str)
		if err != nil {
			panic(err)
		}
		//defer session.Close()

		// Optional. Switch the session to a monotonic behavior.
		//session.SetMode(mgo.Monotonic, true)

		sess, _ := c.Otto.ToValue(session)
		return sess
	})

	obj.Set("one", func(c otto.FunctionCall) otto.Value {
		queryObj, _ := c.Argument(0).Export()
		query, ok := queryObj.(*mgo.Query)
		if ok {
			out := map[string]interface{}{}
			query.One(&out)
			val, _ := vm.ToValue(out)
			return val
		}
		return otto.FalseValue()
	})

	obj.Set("all", func(c otto.FunctionCall) otto.Value {
		queryObj, _ := c.Argument(0).Export()
		query, ok := queryObj.(*mgo.Query)
		if ok {
			out := []map[string]interface{}{}
			query.All(&out)
			val, _ := vm.ToValue(out)
			return val
		}
		return otto.FalseValue()
	})

	return obj.Value()
}

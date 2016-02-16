package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"./../modules"
	"github.com/robertkrimen/otto"
	"github.com/satori/go.uuid"
)

func InitPlugin() *modules.Plugin {
	p := modules.Plugin{
		Name: "crypto",
		Init: registerVM,
	}
	return &p
}

func registerVM(vm *otto.Otto) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("md5", func(c otto.FunctionCall) otto.Value {
		arg1, _ := c.Argument(0).ToString()
		k := md5.Sum([]byte(arg1))
		resp := hex.EncodeToString(k[:])
		oobj, _ := otto.ToValue(resp)
		return oobj
	})
	obj.Set("sha256", func(c otto.FunctionCall) otto.Value {
		arg1, _ := c.Argument(0).ToString()
		k := sha256.Sum256([]byte(arg1))
		resp := hex.EncodeToString(k[:])
		oobj, _ := otto.ToValue(resp)
		return oobj
	})
	obj.Set("newGuid", func(c otto.FunctionCall) otto.Value {
		u1 := uuid.NewV4()
		oobj, _ := otto.ToValue(u1.String())
		return oobj
	})
	return obj.Value()
}

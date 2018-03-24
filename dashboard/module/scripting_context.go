package module

import (
	"github.com/yuin/gopher-lua"
	"net/url"
)

const (
	luaModuleTypeName = "module"
	luaModuleTypeValues = "urlValues"
)

// Registers my person type to given L.
func registerModuleType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaModuleTypeName)
	L.SetGlobal("module", mt)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), contextMethods))
}

// Checks whether the first lua argument is a *LUserData with *Module and returns this *Module.
func checkModule(L *lua.LState) *Module {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Module); ok {
		return v
	}
	L.ArgError(1, "module expected")
	return nil
}

var contextMethods = map[string]lua.LGFunction{
	"render": moduleRender,
}

// Getter and setter for the Person#Name
func moduleRender(L *lua.LState) int {
	c := checkModule(L)

	args := L.GetTop()

	if args >= 2 {
		var data interface{}

		view := L.CheckString(2)

		if args == 3 {
			v := L.CheckTable(3)

			data = ToGoValue(v)
		}

		L.Push(lua.LString(c.Render(view, data)))
		return 1
	}
	L.ArgError(1, "expecting 2 or more arguments")
	return 0
}

func registerValuesType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaModuleTypeValues)
	L.SetGlobal("urlValues", mt)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), valuesMethods))
}

func checkValues(L *lua.LState) url.Values {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(url.Values); ok {
		return v
	}
	L.ArgError(1, "url.Values expected")
	return nil
}

var valuesMethods = map[string]lua.LGFunction{
	"get": valuesGet,
}

func valuesGet(L *lua.LState) int {
	v := checkValues(L)

	if L.GetTop() == 2 {
		L.Push(lua.LString(v.Get(L.CheckString(2))))
		return 1
	}

	L.ArgError(1, "expecting 1 or more arguments")
	return 0
}
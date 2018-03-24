package module

import (
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	luajson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"

	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func (m *Module) initializeScripts() {
	if m.FileExists(m.ID + ".lua") {
		L := lua.NewState()

		L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
		L.PreloadModule("json", luajson.Loader)
		L.PreloadModule("url", gluaurl.Loader)

		L.DoFile(m.FilePath(m.ID + ".lua"))

		registerModuleType(L)
		registerValuesType(L)

		ud := L.NewUserData()
		ud.Value = m
		L.SetMetatable(ud, L.GetTypeMetatable(luaModuleTypeName))

		L.SetGlobal("ctx", ud)

		if m.Config != nil {
			L.SetGlobal("config", ToLuaValue(m.Config))
		} else {
			L.SetGlobal("config", &lua.LTable{})
		}

		m.LuaState = L

		m.CallHandler("initialize", nil)
	}
}

func (m *Module) CallHandler(handler string, vars url.Values) (interface{}, error) {
	v := m.LuaState.GetGlobal(handler)

	if v == nil || v == lua.LNil {
		return nil, errors.New("unknown handler " + handler)
	}

	param := lua.P{
		Fn:      v,
		NRet:    1,
		Protect: true,
	}

	vals := m.LuaState.NewUserData()
	vals.Value = vars
	m.LuaState.SetMetatable(vals, m.LuaState.GetTypeMetatable(luaModuleTypeValues))

	if err := m.LuaState.CallByParam(param, vals); err != nil {
		return nil, err
	}

	ret := m.LuaState.Get(-1)

	m.LuaState.Pop(1)

	return ToGoValue(ret), nil
}

func ToGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key))
				ret[keystr] = ToGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}

func ToLuaValue(gv interface{}) lua.LValue {
	switch v := gv.(type) {
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case int:
	case int8:
	case uint8:
	case int16:
	case uint16:
	case int32:
	case uint32:
	case int64:
	case uint64:
	case float32:
	case float64:
		return lua.LNumber(float64(v))
	case map[string]interface{}:
		t := &lua.LTable{}

		for k, val := range v {
			t.RawSetString(k, ToLuaValue(val))
		}

		return t
	}
	return lua.LNil
}
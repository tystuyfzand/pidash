package module

import (
	"testing"
	"github.com/go-ini/ini"
	"github.com/yuin/gopher-lua"
)

func TestModule_Config(t *testing.T) {
	cfg := []byte("test = \"test\"")

	f, err := ini.Load(cfg)

	if err != nil {
		t.Fatal(err)
	}

	m := mapConfig(f.Section(""))

	L := lua.NewState()

	L.SetGlobal("config", ToLuaValue(m))

	if err := L.DoString("function getTest() return config.test end"); err != nil {
		t.Fatal(err)
	}

	fn := L.GetGlobal("getTest").(*lua.LFunction)

	err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	ret := L.Get(-1)

	L.Pop(1)

	if ret.String() != "test" {
		t.Fatal("Expected value test, got", ret.String())
	}
}
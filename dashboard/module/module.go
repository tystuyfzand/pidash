package module

import (
	"path"
	"os"
	"encoding/json"
	"github.com/yuin/gopher-lua"
	"html/template"
	"log"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"github.com/tystuyfzand/pidash/dashboard/october"
)

var (
	modules = make(map[string]*Module)
)

type Module struct {
	ID string         `json:"-"`
	Config map[string]interface{} `json:"-"`
	Name string       `json:"name"`
	Assets *Assets    `json:"assets"`
	DefaultSize []int `json:"defaultSize"`

	Templates map[string]*template.Template `json:"-"`
	LuaState *lua.LState `json:"-"`
}

type Assets struct {
	Views []string `json:"views"`
	Styles []string `json:"styles"`
	Scripts []string `json:"scripts"`
}

func (m *Module) Initialize() {
	if config.Debug {
		log.Println("Initializing module", m.ID)
	}

	m.initializeConfig()
	m.initializeAssets()
	m.initializeViews()
	m.initializeScripts()
}

func (m *Module) FilePath(file string) string {
	return path.Join(config.Config.Modules.Path, m.ID, file)
}

func (m *Module) FileExists(path string) bool {
	if _, err := os.Stat(m.FilePath(path)); os.IsNotExist(err) {
		return false
	}

	return true
}

func (m *Module) OctoberAssets() *october.OctoberAssets {
	var assets *october.OctoberAssets

	if m.Assets != nil {
		assets = &october.OctoberAssets{}

		if len(m.Assets.Styles) > 0 {
			assets.Styles = make([]string, len(m.Assets.Styles))

			for i, style := range m.Assets.Styles {
				assets.Styles[i] = style
			}
		}

		if len(m.Assets.Scripts) > 0 {
			assets.Scripts = make([]string, len(m.Assets.Scripts))

			for i, script := range m.Assets.Scripts {
				assets.Scripts[i] = script
			}
		}
	}

	return assets
}

func Load(id string) (*Module, error) {
	definitionPath := path.Join(config.Config.Modules.Path, id, "module.json")

	if _, err := os.Stat(definitionPath); os.IsNotExist(err) {
		return nil, err
	}

	if config.Debug {
		log.Println("Loading module", id)
	}

	f, err := os.Open(definitionPath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	module := new(Module)

	module.ID = id

	if err := json.NewDecoder(f).Decode(module); err != nil {
		return nil, err
	}

	module.Initialize()

	modules[module.ID] = module

	return module, nil
}

func Get(id string) *Module {
	return modules[id]
}

func List() map[string]*Module {
	return modules
}
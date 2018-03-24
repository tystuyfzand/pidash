package module

import (
	"github.com/tystuyfzand/pidash/dashboard/config"
	"path"
	"os"
	"github.com/go-ini/ini"
)

func (m *Module) initializeConfig() {
	preferredPath := path.Join(config.Config.Modules.ConfigPath, m.Name + ".conf")

	if _, err := os.Stat(preferredPath); !os.IsNotExist(err) {
		// Load preferred
		m.Config = loadConfigFrom(preferredPath)

		if m.Config != nil {
			return
		}
	}

	secondaryPath := path.Join(config.Config.Modules.Path, m.ID, "module.conf")

	if _, err := os.Stat(secondaryPath); !os.IsNotExist(err) {
		// Load included module config
		m.Config = loadConfigFrom(secondaryPath)
	}
}

func loadConfigFrom(file string) map[string]interface{} {
	f, err := ini.Load(file)

	if err != nil {
		return nil
	}

	section, err := f.GetSection("")

	if err != nil {
		return nil
	}

	return mapConfig(section)
}

func mapConfig(section *ini.Section) map[string]interface{} {
	table := make(map[string]interface{})

	for _, v := range section.Keys() {
		if fv, err := v.Float64(); err == nil {
			table[v.Name()] = fv
		} else if bv, err := v.Bool(); err == nil {
			table[v.Name()] = bv
		} else {
			table[v.Name()] = v.String()
		}
	}

	return table
}
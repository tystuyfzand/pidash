package module

import (
	"github.com/tystuyfzand/pidash/dashboard/assets"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"log"
	"path"
)

func (m *Module) initializeAssets() {
	if m.Assets == nil {
		return
	}

	if config.Debug {
		log.Println("Compiling assets for", m.ID)
	}

	if err := initializeCssAssets(m); err != nil {
		log.Fatalln("Unable to compile css assets for", m.ID, ":", err)
	}

	if err := initializeScriptAssets(m); err != nil {
		log.Fatalln("Unable to compile js assets for", m.ID, ":", err)
	}
}

func initializeCssAssets(m *Module) error {
	styles := make([]string, len(m.Assets.Styles))

	for i, style := range m.Assets.Styles {
		styles[i] = m.FilePath("css/" + style)
	}

	if err := assets.CompileCSS(styles, m.ID + ".min.css"); err == nil {
		m.Assets.Styles = []string{"assets/css/compiled/" +  m.ID + ".min.css"}
	} else {
		moduleStyles := make([]string, len(m.Assets.Styles))

		for i, style := range m.Assets.Styles {
			moduleStyles[i] = path.Join("modules", m.ID, "css", style)
		}

		m.Assets.Styles = moduleStyles
	}

	return nil
}

func initializeScriptAssets(m *Module) error {
	scripts := make([]string, len(m.Assets.Scripts))

	for i, script := range m.Assets.Scripts {
		scripts[i] = m.FilePath("js/" + script)
	}

	assets.CompileJS(scripts, m.ID + ".min.js")

	m.Assets.Scripts = []string{"assets/js/compiled/" +  m.ID + ".min.js"}

	return nil
}
package dashboard

import (
	"github.com/tystuyfzand/pidash/dashboard/module"
	"github.com/tystuyfzand/pidash/dashboard/assets"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"log"
	"io/ioutil"
	"path"
)

func Dashboard() {
	compileApplicationAssets()

	loadModules()
}

func compileApplicationAssets() {
	if config.Debug {
		log.Println("Compiling application assets")
	}

	styles := []string{
		"html/assets/css/gridstack.min.css",
		"html/assets/css/fa-solid.css",
		"html/assets/css/fontawesome.css",
		"html/assets/css/weather-icons.css",
		"html/assets/css/bootstrap.css",
		"html/assets/css/style.css",
	}

	for i, style := range styles {
		styles[i] = path.Join(config.DataDirectory, style)
	}

	err := assets.CompileCSS(styles, "app.min.css")

	if err != nil {
		log.Fatalln("Unable to compile css:", err)
	}

	scripts := []string{
		"html/assets/js/jquery-3.3.1.min.js",
		"html/assets/js/jquery-ui.min.js",
		"html/assets/js/bootstrap.js",
		"html/assets/js/bootbox.js",
		"html/assets/js/mousetrap.js",
		"html/assets/js/notify.js",
		"html/assets/js/october/framework.js",
		"html/assets/js/october/framework.extras.js",
		"html/assets/js/october/assetmanager.js",
		"html/assets/js/lodash.min.js",
		"html/assets/js/gridstack.js",
		"html/assets/js/gridstack.jQueryUI.js",
		"html/assets/js/dashboard.js",
	}

	for i, script := range scripts {
		scripts[i] = path.Join(config.DataDirectory, script)
	}

	err = assets.CompileJS(scripts, "app.min.js")

	if err != nil {
		log.Fatalln("Unable to compile js:", err)
	}
}

func loadModules() {
	files, err := ioutil.ReadDir(config.Config.Modules.Path)

	if err != nil {
		log.Fatalln(err)
	}

	var name string

	for _, file := range files {
		name = file.Name()

		if name[0] == '.' || !file.IsDir() {
			continue
		}

		if _, err := module.Load(name); err != nil {
			log.Fatalln("Unable to load module " + name, err)
		}
	}
}
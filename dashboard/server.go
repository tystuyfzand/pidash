package dashboard

import (
	"github.com/tystuyfzand/pidash/dashboard/module"
	"github.com/tystuyfzand/pidash/dashboard/october"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"regexp"
	"net/http"
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
	"io"
	"path"
	"fmt"
	"strings"
)

var (
	cssFile = regexp.MustCompile("\\.css$")
	jsFile = regexp.MustCompile("\\.js$")

	moduleRegexp = regexp.MustCompile("(.*?)::(.*)")

	fileServer http.Handler
	moduleServer http.Handler
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	if handler := r.Header.Get("X-October-Request-Handler"); handler != "" {
		processAjaxHandler(w, r, handler)
		return
	}

	uri := r.RequestURI

	if cssFile.MatchString(uri) {
		w.Header().Set("Content-Type", "text/css")
	} else if jsFile.MatchString(uri) {
		w.Header().Set("Content-Type", "application/javascript")
	} else {
		w.Header().Set("Content-Type", "text/html")
	}

	// Serve LESS/CSS directly from the module path
	if strings.Index(uri, "/modules") == 0 && cssFile.MatchString(uri){
		r.RequestURI = uri[9:]

		moduleServer.ServeHTTP(w, r)
		return
	}

	fileServer.ServeHTTP(w, r)
}

func processAjaxHandler(w http.ResponseWriter, r *http.Request, handler string) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	if handler == "onLoadLayout" {
		f, err := os.Open(path.Join(config.DataDirectory, "modules/config.json"))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		defer f.Close()

		_, err = io.Copy(w, f)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else if handler == "onSaveLayout" {
		// Save data
		layout := r.Form.Get("layout")

		if layout != "" {
			err := ioutil.WriteFile(path.Join(config.DataDirectory, "modules/config.json"), []byte(layout), 0644)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	} else if handler == "onListWidgets" {
		json.NewEncoder(w).Encode(module.List())
	} else if m := moduleRegexp.FindStringSubmatch(handler); m != nil {
		if mp := module.Get(m[1]); mp != nil {
			t, err := mp.CallHandler(m[2], r.Form)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			var res map[string]interface{}

			switch t.(type) {
			case map[string]interface{}:
				res = t.(map[string]interface{})
			case string:
				res = map[string]interface{}{
					"#" + r.Form.Get("id") + " > .grid-stack-item-content" : t.(string),
				}
			default:
				res = make(map[string]interface{})
			}

			if m[2] == "onRender" {
				res[october.Assets] = mp.OctoberAssets()
			}

			if err = json.NewEncoder(w).Encode(res); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No module registered for " + m[1]))
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("No handler for " + handler))
	}
}

func Serve() {
	fileServer = http.FileServer(http.Dir(path.Join(config.DataDirectory, "html")))
	moduleServer = http.FileServer(http.Dir(path.Join(config.Config.Modules.Path)))

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)

	listenAddress := fmt.Sprintf("%s:%d", config.Config.Http.Host, config.Config.Http.Port)

	log.Println("Listening on " + listenAddress)

	http.ListenAndServe(listenAddress, mux)
}

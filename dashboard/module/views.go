package module

import (
	"path"
	"io/ioutil"
	"bytes"
	"html/template"
	"log"
)

func (m *Module) initializeViews() error {
	if m.Assets == nil || len(m.Assets.Views) == 0 {
		return nil
	}

	m.Templates = make(map[string]*template.Template)

	var ext, viewName string

	for _, view := range m.Assets.Views {
		ext = path.Ext(view)
		viewName = view[0:len(view)-len(ext)]

		tpl := template.New(viewName)

		b, err := ioutil.ReadFile(m.FilePath("views/" + view))

		if err != nil {
			return err
		}

		if tpl, err = tpl.Parse(string(b)); err != nil {
			return err
		}

		m.Templates[viewName] = tpl
	}

	return nil
}

func (m *Module) Render(view string, data interface{}) string {
	tpl, exists := m.Templates[view]

	if !exists {
		return ""
	}

	var b bytes.Buffer

	if err := tpl.Execute(&b, data); err != nil {
		log.Println("Unable to render", view, ":", err)
	}

	return b.String()
}

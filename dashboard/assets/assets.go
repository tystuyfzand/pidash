package assets

import (
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"os"
	"io/ioutil"
	"path"
	"bytes"
	"regexp"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"os/exec"
	"errors"
)

var (
	lessFile = regexp.MustCompile("\\.less$")
)

func CompileCSS(files []string, outputName string) error {
	var buf bytes.Buffer

	for _, file := range files {
		if lessFile.MatchString(file) {
			if !config.HasLessCompiler {
				return errors.New("no less compiler found")
			}

			outputPath := path.Join(os.TempDir(), path.Base(lessFile.ReplaceAllString(file, ".css")))

			defer func(p string) {
				// Some kind of sanity checking, I guess.
				if s, err := os.Stat(p); s.IsDir() || os.IsNotExist(err) {
					return
				}

				os.Remove(p)
			}(outputPath)

			lessToCss(file, outputPath)

			file = outputPath
		}

		b, err := ioutil.ReadFile(file)

		if err != nil {
			return err
		}

		buf.Write(b)
	}

	minifiedString, err := makeMinified(buf.String(), "text/css")

	if err != nil {
		return err
	}

	basePath := path.Join(config.DataDirectory, "html/assets/css/compiled")

	os.MkdirAll(basePath, 0755)

	return ioutil.WriteFile(path.Join(basePath, outputName), []byte(minifiedString), 0644)
}

func CompileJS(files []string, outputName string) error {
	var buf bytes.Buffer

	for _, file := range files {
		b, err := ioutil.ReadFile(file)

		if err != nil {
			return err
		}

		buf.Write(b)
	}

	minifiedString, err := makeMinified(buf.String(), "text/javascript")

	if err != nil {
		return err
	}

	basePath := path.Join(config.DataDirectory, "html/assets/js/compiled")

	os.MkdirAll(basePath, 0755)

	return ioutil.WriteFile(path.Join(basePath, outputName), []byte(minifiedString), 0644)
}

func makeMinified(input, fileType string) (string, error) {
	mini := minify.New()

	mini.AddFunc("text/css", css.Minify)
	mini.AddFunc("text/javascript", js.Minify)

	return mini.String(fileType, input)
}

func lessToCss(input, output string) error {
	cmd := exec.Command("lessc", input, output)

	return cmd.Run()
}

func HasLessCompiler() bool {
	cmd := exec.Command("lessc", "--version")

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
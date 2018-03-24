package config

import (
	"os"
	"github.com/go-ini/ini"
	"path/filepath"
)

var (
	Config *Dashboard

	Debug bool
	HasLessCompiler bool

	DataDirectory string
)

type Dashboard struct {
	Http HttpStruct `ini:"http"`
	Modules ModuleStruct `ini:"modules"`
}

type HttpStruct struct {
	Host string `ini:"host"`
	Port int `ini:"port"`
}

type ModuleStruct struct {
	Path string `ini:"path"`
	ConfigPath string `ini:"configPath"`
}

func Load(config string) error {
	f, err := os.Open(config)

	if err != nil {
		return err
	}

	defer f.Close()

	Config = new(Dashboard)

	if err := ini.MapTo(Config, f); err != nil {
		return err
	}

	checkDefaults(config)

	return nil
}

func checkDefaults(config string) {
	if Config.Modules.Path == "" {
		Config.Modules.Path = filepath.Join(DataDirectory, "modules")
	}

	if Config.Modules.ConfigPath == "" {
		Config.Modules.ConfigPath = filepath.Dir(config)
	} else if !filepath.IsAbs(Config.Modules.ConfigPath) {
		Config.Modules.ConfigPath, _ = filepath.Abs(filepath.Join(filepath.Dir(config), Config.Modules.ConfigPath))
	}
}
package gofigure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigType string

const (
	GlobalConfig ConfigType = "GlobalConfig"
	CacheConfig  ConfigType = "CacheConfig"
	LocalConfig  ConfigType = "LocalConfig"
)

const (
	CONFIG_FILE_NAME = "config.yml"
	CACHE_FILE_NAME  = "cache.yml"
)

type Configuration struct {
	Entity    any
	Path      string
	Directory string
	Error     error
}

type ConfigurationsType map[ConfigType]*Configuration

var (
	Project        string
	Configurations = ConfigurationsType{}
)

var (
	ErrorConfigNotRegistered = func(typ ConfigType) error { return errors.New(fmt.Sprintf("%s not registered.", typ)) }
)

func CheckErr(typ ConfigType) {
	var (
		err error
		ok  bool
		cfg *Configuration
	)
	var handleError = func(err error) {
		fmt.Sprintf("%s[%s]%s %s\n", "\033[31m", "ERROR", "\033[0m", err.Error())
		os.Exit(1)
	}
	if cfg, ok = Configurations[typ]; !ok {
		handleError(ErrorConfigNotRegistered(typ))
	}
	err = cfg.Error
	if err != nil {
		handleError(err)
	}
}

func safeMakeFile(dir string, base string) (string, error) {
	var (
		err       error
		osFilePth string
	)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	osFilePth = path.Join(dir, base)

	if _, err = os.Stat(osFilePth); os.IsNotExist(err) {
		fl, err := os.Create(osFilePth)
		if err != nil {
			return "", err
		}
		err = fl.Close()
		if err != nil {
			return "", err
		}
	}
	return osFilePth, nil
}

func initConfig(typ ConfigType) (string, string, error) {
	var (
		err                            error
		osPth, cfgFlNm, cfgDir, cfgPth string
	)

	switch typ {
	case GlobalConfig:
		osPth, err = os.UserConfigDir()
		cfgFlNm = CONFIG_FILE_NAME
	case CacheConfig:
		osPth, err = os.UserCacheDir()
		cfgFlNm = CACHE_FILE_NAME
	case LocalConfig:
		// TODO: populate local config
		cfgFlNm = CONFIG_FILE_NAME
		return "", "", nil
	}
	cfgDir = filepath.Join(osPth, Project)

	if err != nil {
		return "", cfgDir, nil
	}

	cfgPth, err = safeMakeFile(
		cfgDir,
		cfgFlNm,
	)
	return cfgPth, cfgDir, err
}

func hydrate(cfg *Configuration) error {
	var (
		err      error
		fileByte []byte
	)

	fileByte, err = ioutil.ReadFile(cfg.Path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(fileByte, cfg.Entity)
}

func Setup(prj string) {
	Project = prj
}

func Register[T any](typ ConfigType, ent *T) {
	var (
		err error
		cfg = &Configuration{Entity: ent}
	)
	cfg.Path, cfg.Directory, err = initConfig(typ)
	if err != nil {
		cfg.Error = err
	}
	err = hydrate(cfg)
	if err != nil {
		cfg.Error = err
	}
	Configurations[typ] = cfg
}

func WriteAll() {
	var (
		err    error
		cfg    *Configuration
		typ    ConfigType
		cfgByt []byte
	)

	for _, cfg = range Configurations {
		if cfg.Error != nil {
			continue
		}
		cfgByt, err = yaml.Marshal(&cfg.Entity)
		if err != nil {
			cfg.Error = err
			continue
		}

		err = ioutil.WriteFile(cfg.Path, cfgByt, 0777)
		if err != nil {
			cfg.Error = err
			continue
		}
	}
	for typ = range Configurations {
		CheckErr(typ)
	}
}

func RemoveAll() error {
	var (
		err error
		cfg *Configuration
	)

	for _, cfg = range Configurations {
		err = os.RemoveAll(cfg.Directory)
		if err != nil {
			return err
		}
	}
	return nil
}

func Where() {
	var (
		cfg *Configuration
		typ ConfigType
	)
	for typ, cfg = range Configurations {
		fmt.Println("[", typ, "]", cfg.Path)
	}
}

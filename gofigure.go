package gofigure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/nanvenomous/exfs"
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

type initFunc func() (string, string, error)

type Configuration struct {
	Entity    any
	Path      string
	Directory string
	Error     error
}

type ConfigRegistration[T any] struct {
	Type          ConfigType
	InitialConfig *T
}

type ConfigurationsType map[ConfigType]*Configuration

var (
	fs             = exfs.NewFileSystem()
	Project        string
	Configurations = ConfigurationsType{}
)

var (
	// ErrorGeneral             = func(typ ConfigType, err error) error { return errors.New(fmt.Sprintf("%s %s", typ, err.Error())) }
	ErrorConfigNotRegistered = func(typ ConfigType) error { return errors.New(fmt.Sprintf("%s not registered.", typ)) }
	ErrorLocateConfig        = func(typ ConfigType, msg string) error {
		return errors.New(fmt.Sprintf("Could not locate %s. Are you in the correct directory?\n%s", typ, msg))
	}
)

func localProjectFolder() string {
	return "." + Project
}

func handleError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s[%s]%s %s", "\033[31m", "ERROR", "\033[0m", err.Error()))
		os.Exit(1)
	}
}

func GetErr(typ ConfigType) error {
	var (
		ok  bool
		cfg *Configuration
	)
	if cfg, ok = Configurations[typ]; !ok {
		return ErrorConfigNotRegistered(typ)
	}
	return cfg.Error
}

func CheckErr(typ ConfigType) {
	handleError(GetErr(typ))
}

func CheckAllErrs() {
	var typ ConfigType
	for typ = range Configurations {
		CheckErr(typ)
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
		osPthFunc                      func() (string, error)
		osPth, cfgFlNm, cfgDir, cfgPth string
	)

	if typ == CacheConfig {
		osPthFunc = os.UserCacheDir
		cfgFlNm = CACHE_FILE_NAME
	} else {
		osPthFunc = os.UserConfigDir
		cfgFlNm = CONFIG_FILE_NAME
	}

	switch typ {
	case GlobalConfig, CacheConfig:
		osPth, err = osPthFunc()
		if err != nil {
			return "", cfgDir, err
		}
		cfgDir = filepath.Join(osPth, Project)
	case LocalConfig:
		cfgDir, err = fs.FindFileInAboveCurDir(localProjectFolder())
		if err != nil {
			return "", cfgDir, ErrorLocateConfig(LocalConfig, err.Error())
		}
		return filepath.Join(cfgDir, cfgFlNm), cfgDir, nil
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

func Register[T any](typ ConfigType, initialCfg *T) {
	var curInitFunc = func() (string, string, error) {
		return initConfig(typ)
	}
	Configurations[typ] = register(curInitFunc, initialCfg)
}

func register[T any](intf initFunc, ent *T) *Configuration {
	var (
		err error
		cfg = &Configuration{Entity: ent}
	)
	cfg.Path, cfg.Directory, err = intf()
	if err != nil {
		cfg.Error = err
		return cfg
	}
	err = hydrate(cfg)
	if err != nil {
		cfg.Error = err
	}
	return cfg
}

func Setup(prj string) {
	Project = prj
}

func CreateLocalConfig[T any](initialLocCfg *T) error {
	var createLocCfg = func() (string, string, error) {
		var (
			err                error
			wd, cfgDir, cfgPth string
		)
		wd, err = os.Getwd()
		if err != nil {
			return "", "", err
		}

		cfgDir = filepath.Join(wd, localProjectFolder())

		cfgPth, err = safeMakeFile(
			cfgDir,
			CONFIG_FILE_NAME,
		)
		return cfgPth, cfgDir, nil
	}

	Configurations[LocalConfig] = register(createLocCfg, initialLocCfg)
	return Configurations[LocalConfig].Error
}

func WriteAll() {
	var (
		err    error
		cfg    *Configuration
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

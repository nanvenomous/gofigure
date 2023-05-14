package gofigure

import (
	"os"
	"testing"
)

const PROJECT_NAME = "my-gofigure-test-project"

// Define the structure of your config
type myGlobalConfigType struct {
	User struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	} `yaml:"user"`
	VerbosityLevel uint8 `yaml:"verbosity_level"`
}

// Create a config entity, set default values (if desired)
var myGlobalConfig = &myGlobalConfigType{VerbosityLevel: 1}

// Define the explicit interface to your configuration
func SetUsername(unm string) {
	CheckErr(GlobalConfig)
	myGlobalConfig.User.Name = unm
}
func Username() string {
	CheckErr(GlobalConfig)
	return myGlobalConfig.User.Name
}
func VerbosityLevel() uint8 {
	CheckErr(GlobalConfig)
	return myGlobalConfig.VerbosityLevel
}

func setup(m *testing.M) error {
	var (
		err error
	)
	Setup(PROJECT_NAME)
	Register(GlobalConfig, myGlobalConfig)

	err = RemoveAll()
	if err != nil {
		return err
	}

	code := m.Run()

	err = RemoveAll()
	if err != nil {
		return err
	}

	os.Exit(code)
	return nil
}

func TestMain(m *testing.M) {
	err := setup(m)
	if err != nil {
		panic(err)
	}
}

// Create & update your config
func TestGlobalConfigCreation(t *testing.T) {
	Setup(PROJECT_NAME)
	Register(GlobalConfig, myGlobalConfig)

	if usr := Username(); usr != "" {
		t.Error("username was set, i.e. accidentially persisted")
	}

	if vl := VerbosityLevel(); vl != 1 {
		t.Error("verbosity was not set as default")
	}

	SetUsername("nanvenomous")

	if usr := Username(); usr != "nanvenomous" {
		t.Error("failed to set username")
	}

	// write the config to disk
	WriteAll()
}

// Retrieve config from disk
func TestGlobalConfigRetrieval(t *testing.T) {
	Setup(PROJECT_NAME)
	Register(GlobalConfig, myGlobalConfig)

	// empty the configuration entity
	myGlobalConfig = &myGlobalConfigType{}
	Register(GlobalConfig, myGlobalConfig)

	if usr := Username(); usr != "nanvenomous" {
		t.Error("the username was not persisted")
	}

	if vl := VerbosityLevel(); vl != 1 {
		t.Error("the verbosity was not persisted")
	}

}

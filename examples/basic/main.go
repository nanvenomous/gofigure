package main

import (
	"fmt"

	"github.com/nanvenomous/gofigure"
)

// name your project (to be used for directory naming)
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
	gofigure.CheckErr(gofigure.GlobalConfig) // check for file/parsing errors only when accessing config, or replace with your own check, you do you, boo
	myGlobalConfig.User.Name = unm
}
func Username() string {
	gofigure.CheckErr(gofigure.GlobalConfig)
	return myGlobalConfig.User.Name
}
func VerbosityLevel() uint8 {
	gofigure.CheckErr(gofigure.GlobalConfig)
	return myGlobalConfig.VerbosityLevel
}

func main() {
	var err error
	// initialize files, set defaults
	gofigure.Setup(PROJECT_NAME)
	gofigure.Register(gofigure.GlobalConfig, myGlobalConfig) // gofigure supports Global, Local, and Cache configurations

	// access the config data
	fmt.Println(Username())       // ""
	SetUsername("nanvenomous")    //
	fmt.Println(Username())       // "nanvenomous"
	fmt.Println(VerbosityLevel()) // 1

	// write the config to disk, so you can get it on next program run
	gofigure.WriteAll() // uses go standard lib to coose location by your project name and OS

	// prints the paths to all your registered config files (in case your user is wondering :)
	gofigure.Where() // output depends on your os & your project name

	// remove all configuration files (really nice for uninstall scripts)
	if err = gofigure.RemoveAll(); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"

	"github.com/nanvenomous/gofigure"
)

// name your project (to be used for directory naming)
const PROJECT_NAME = "my-gofigure-test-project"

// Define the structure of your config
type myLocalConfigType struct {
	Bare bool `yaml:"bare"`
}
type myCacheConfigType struct {
	History []string `yaml:"history"`
}

// Create Local Config
var myLocalConfig = &myLocalConfigType{Bare: false}

func SetBare(br bool) {
	gofigure.CheckErr(gofigure.LocalConfig)
	myLocalConfig.Bare = br
}
func BareE() (bool, error) { // return the error in case you want to continue project execution, or maybe replace with your own error
	var err = gofigure.GetErr(gofigure.LocalConfig)
	if err != nil {
		return false, err
	}
	return myLocalConfig.Bare, nil
}

// Create Cache Config
var myCacheConfig = &myCacheConfigType{History: []string{}}

func SetHistory(hst []string) {
	gofigure.CheckErr(gofigure.CacheConfig)
	myCacheConfig.History = hst
}
func History() []string {
	gofigure.CheckErr(gofigure.CacheConfig)
	return myCacheConfig.History
}

func main() {
	var err error
	// set project name, register local config type
	gofigure.Setup(PROJECT_NAME)
	gofigure.Register(gofigure.LocalConfig, myLocalConfig)
	gofigure.Register(gofigure.CacheConfig, myCacheConfig)

	// access the config data
	_, err = BareE()
	fmt.Println(err) // "Could not locate LocalConfig. Are you in the correct directory?"
	// local config needs to be created. (for example the 'git init' command)
	err = gofigure.CreateLocalConfig(myLocalConfig)
	if err != nil {
		fmt.Println(err) // nil
	}
	// now we have access to the local config from anywhere under our current directory
	isBare, err := BareE()
	if err != nil {
		fmt.Println(err) // nil
	}
	fmt.Println(isBare) // false

	SetHistory([]string{"git status", "cd ..", "whoami"})
	SetHistory(append(History()[1:], "timedatectl"))
	fmt.Println(History()) // "[cd .. whoami timedatectl]"

	gofigure.WriteAll()
	gofigure.Where() // we can see the local config was created in our current directory
	if err = gofigure.RemoveAll(); err != nil {
		panic(err)
	}
}

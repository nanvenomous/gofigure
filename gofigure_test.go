package gofigure

import "testing"

// Define the structure of your config
type myGlobalConfigType struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
}

// Assign default values
var myGlobalConfig = &myGlobalConfigType{}

// Define the explicit interface to your configuration
func SetUsername(unm string) {
	CheckErr(GlobalConfig)
	myGlobalConfig.Username = unm
}
func Username() string {
	CheckErr(GlobalConfig)
	return myGlobalConfig.Username
}

// Access & edit your config
func TestGlobalConfig(t *testing.T) {
	Setup("my-project")
	Register(GlobalConfig, myGlobalConfig)

	if usr := Username(); usr != "" {
		t.Error("username is initiall an empty string")
	}

	SetUsername("nanvenomous")

	if usr := Username(); usr != "nanvenomous" {
		t.Error("failed to set username")
	}
}

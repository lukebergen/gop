package gop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	Config  Configuration
	GopDir  string = filepath.Join(os.Getenv("HOME"), ".gopivot")
	DbDir   string = filepath.Join(GopDir, "database")
	Version string
)

type Configuration struct {
	CurrentUser      User
	CurrentProjectId int
}

type Completion struct {
	Id           int
	Text         string
	CurrentState string
	LastTouched  time.Time
}

func Init(version string) {
	Version = version
	if err := os.MkdirAll(DbDir, 0700); err != nil {
		fmt.Println("AWWW SHIT!")
		panic(err)
	}
	LoadConfig()
}

func LoadConfig() {
	if err := os.MkdirAll(GopDir, 0700); err != nil {
		panic(err)
	}

	configFile, err := os.OpenFile(filepath.Join(GopDir, "config.json"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	json.NewDecoder(configFile).Decode(&Config)
}

func SaveConfig() {
	fileJson, err := json.Marshal(Config)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(filepath.Join(GopDir, "config.json"), fileJson, 0600); err != nil {
		panic(err)
	}
}

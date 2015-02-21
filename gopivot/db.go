package gop

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	Config Configuration
	GopDir string = filepath.Join(os.Getenv("HOME"), ".gopivot")
	DbDir  string = filepath.Join(GopDir, "database")
)

type Configuration struct {
	CurrentUser      User
	CurrentProjectId int
}

func Init() {
	if err := os.MkdirAll(DbDir, 0744); err != nil {
		panic(err)
	}
	LoadConfig()
}

func LoadConfig() {
	if err := os.MkdirAll(GopDir, 0744); err != nil {
		panic(err)
	}

	configFile, err := os.OpenFile(GopDir+"/config.json", os.O_RDWR|os.O_CREATE, 0640)
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
	if err := ioutil.WriteFile(GopDir+"/config.json", fileJson, 0640); err != nil {
		panic(err)
	}
}

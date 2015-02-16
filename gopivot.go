package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Configuration struct {
	Token string
}

var (
	config Configuration
)

func main() {
	config = GetConfig()
	resp := request("/projects")
	fmt.Printf(resp)
}

func request(endpoint string) string {
	pivotalApiUrl := "https://www.pivotaltracker.com/services/v5"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", pivotalApiUrl+endpoint, nil)
	req.Header.Add("X-TrackerToken", config.Token)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func GetConfig() Configuration {
	_, err := os.Stat("~/.gopivot")
	if os.IsNotExist(err) {
		// get config data
		os.Chdir(os.Getenv("HOME"))
		os.Mkdir("./.gopivot", 0700)
		// write config data
	}
	os.Chdir(os.Getenv("HOME") + "/.gopivot")
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	CurrentUser      User
	CurrentProjectId int
}

type Project struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	Name      string `json:"project_name"`
}

type Story struct {
	AcceptedAt    string   `json:"accepted_at"`
	CurrentState  string   `json:"current_state"`
	Id            int      `json:"id"`
	Labels        []string `json:"labels"`
	Name          string   `json:"name"`
	OwnedById     int      `json:"owned_by_id"`
	OwnerIds      []int    `json:"owner_ids"`
	ProjectId     int      `json:"project_id"`
	RequestedById int      `json:"requested_by_id"`
	StoryType     string   `json:"story_type"`
	Url           string   `json:"url"`
}

type User struct {
	ApiToken string `json:"api_token"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Projects []Project
}

var (
	config    Configuration
	configDir string
)

func main() {
	config = getConfig()
	// simulate various commands

	// login()
	// project("")
	// project("Reviewed")
	// ls()
}

func request(endpoint string) []byte {
	pivotalApiUrl := "https://www.pivotaltracker.com/services/v5"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", pivotalApiUrl+endpoint, nil)
	req.Header.Add("X-TrackerToken", config.CurrentUser.ApiToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes
}

func getConfig() Configuration {
	configDir = os.Getenv("HOME") + "/.gopivot"
	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		os.Mkdir(configDir, 0750)
	}
	_, err = os.Stat(configDir + "/config.json")
	if os.IsNotExist(err) {
		ioutil.WriteFile(configDir+"/config.json", []byte("{}"), 0640)
	}

	jsonBytes, _ := ioutil.ReadFile(configDir + "/config.json")
	var config Configuration
	json.Unmarshal(jsonBytes, &config)
	return config
}

func login() {
	var username, password string
	fmt.Printf("Username: ")
	fmt.Scan(&username)
	fmt.Printf("Password: ")
	fmt.Scan(&password)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.pivotaltracker.com/services/v5/me", nil)
	req.SetBasicAuth(username, password)
	resp, _ := client.Do(req)
	var currentUser User
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &currentUser)

	config = Configuration{CurrentUser: currentUser}
	saveConfig()
}

func project(setTo string) {
	if setTo == "" {
		names := make([]string, len(config.CurrentUser.Projects))
		for i := 0; i < len(names); i++ {
			names[i] = config.CurrentUser.Projects[i].Name
		}
		fmt.Printf(strings.Join(names, ", ") + "\n")
	} else {
		id := -1
		projects := config.CurrentUser.Projects
		for i := 0; i < len(projects); i++ {
			if projects[i].Name == setTo {
				id = projects[i].ProjectId
			}
		}
		if id == -1 {
			fmt.Printf("Not a project name")
		} else {
			config.CurrentProjectId = id
			saveConfig()
		}
	}
}

func ls() {
	stories := make([]Story, 0)
	body := request(fmt.Sprintf("/projects/%v/stories", config.CurrentProjectId))
	json.Unmarshal(body, &stories)

	// TODO: enumerate, store, paginte over, otherwise deal with...
	fmt.Printf("Found %v stories\n", len(stories))
}

func saveConfig() {
	fileJson, _ := json.Marshal(config)
	ioutil.WriteFile(configDir+"/config.json", fileJson, 0640)
}

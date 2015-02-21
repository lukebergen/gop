package gop

import (
	"io/ioutil"
	"net/http"
)

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

func request(endpoint string) ([]byte, http.Header) {
	pivotalApiUrl := "https://www.pivotaltracker.com/services/v5"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", pivotalApiUrl+endpoint, nil)
	req.Header.Add("X-TrackerToken", Config.CurrentUser.ApiToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes, resp.Header
}

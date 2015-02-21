package gop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	flag "github.com/ogier/pflag"
)

type Flags struct {
	Short bool
	User  string
	Debug bool
	State string
}

func Exec() {
	flags := Flags{}

	flag.BoolVarP(&flags.Short, "concise", "c", false, "Print the concise form (similar to git status -s)")
	flag.BoolVarP(&flags.Debug, "debug", "d", false, "Debug mode for dev. Usually just prints the API request being made")
	flag.StringVarP(&flags.User, "user", "u", "me", "Filter by user. 'all' and 'me' are special. Otherwise just a username or initials")
	flag.StringVarP(&flags.State, "state", "s", "active", "comma separated list of states to filter by. Defaults to 'active' which is a gop shorthand for 'started,finished,delivered,rejected'")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, "usage")
	}

	switch args[0] {
	case "login":
		CommandLogin(flags)
	case "project":
		if len(args) == 2 {
			CommandProject(flags, args[1])
		} else {
			CommandProject(flags, "")
		}
	case "ls":
		CommandLs(flags)
	default:
		flag.Usage()
	}
}

func CommandLogin(flags Flags) {
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

	// TODO: use reflection to add Config.Set/Get methods instead of accessing/saving directly like this
	Config.CurrentUser = currentUser
	SaveConfig()
}

func CommandProject(flags Flags, setTo string) {
	if setTo == "" {
		names := make([]string, len(Config.CurrentUser.Projects))
		for i := 0; i < len(names); i++ {
			names[i] = Config.CurrentUser.Projects[i].Name
		}
		fmt.Printf(strings.Join(names, ", ") + "\n")
	} else {
		id := -1
		projects := Config.CurrentUser.Projects
		for i := 0; i < len(projects); i++ {
			if projects[i].Name == setTo {
				id = projects[i].ProjectId
			}
		}
		if id == -1 {
			fmt.Printf("Not a project name")
		} else {
			Config.CurrentProjectId = id
			SaveConfig()
		}
	}
}

func CommandLs(flags Flags) {
	if Config.CurrentProjectId == 0 {
		fmt.Printf("You need to set an active project to list stories\n")
	} else {
		stories := make([]Story, 0)

		qParams := url.Values{}

		var filters []string

		if flags.User != "all" {
			if flags.User == "me" {
				filters = append(filters, "owner:"+Config.CurrentUser.Username)
			} else {
				filters = append(filters, "owner:"+flags.User)
			}
		}

		if flags.State == "active" {
			filters = append(filters, "state:started,finished,delivered,rejected")
		} else {
			filters = append(filters, "state:"+flags.State)
		}

		if len(filters) > 0 {
			qParams.Set("filter", strings.Join(filters, " "))
		}

		reqStr := fmt.Sprintf("/projects/%v/stories?%s", Config.CurrentProjectId, qParams.Encode())
		if flags.Debug {
			fmt.Println("Debug: " + reqStr)
		}

		body, _ := request(reqStr)
		json.Unmarshal(body, &stories)

		if flags.Short {
			for i := 0; i < len(stories); i++ {
				char := strings.ToUpper(string(stories[i].CurrentState[0]))
				fmt.Printf("%v %v: %v\n", char, stories[i].Id, stories[i].Name)
			}
		} else {
			states := []string{"started", "finished", "delivered", "rejected"}
			for i := 0; i < len(states); i++ {
				fmt.Printf("%v:\n", states[i])
				for j := 0; j < len(stories); j++ {
					if stories[j].CurrentState == states[i] {
						fmt.Printf("\t%v: %v\n", stories[j].Id, stories[j].Name)
					}
				}
				fmt.Println()
			}
		}
	}
}

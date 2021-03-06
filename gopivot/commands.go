package gop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/gopass"
	flag "github.com/ogier/pflag"
)

var completionScript string = `
	if [ -n "$BASH_VERSION" ]; then
		_gop_complete() {
			COMPREPLY=()
			local word="${COMP_WORDS[COMP_CWORD]}"
			local completions="$(gop complete "$word")"
			COMPREPLY=( $(compgen -W "$completions" -- "$word") )
		}
		complete -f -F _gop_complete gop
	elif [ -n "$ZSH_VERSION" ]; then
		_gop_complete() {
			local word completions
			word="$1"
			completions="$(gop complete "${word}")"
			reply=( "${(ps:\n:)completions}" )
		}
		compctl -f -K _gop_complete gop
	fi
`

type Flags struct {
	Concise   bool
	Debug     bool
	Help      bool
	ShellInit bool
	State     string
	User      string
	Version   bool // actually handled by main where the const lives
}

type command func(flags Flags)

var commands map[string]command

func Exec() {
	flags := Flags{}

	flag.BoolVarP(&flags.Concise, "concise", "c", false, "Print the concise form (similar to git status -s)")
	flag.BoolVarP(&flags.Debug, "debug", "d", false, "Debug mode for dev. Usually just prints the API request being made")
	flag.BoolVarP(&flags.Help, "help", "", false, "Print help")
	flag.BoolVar(&flags.ShellInit, "shell-init", false, "Generate init shell script. Useful for rc files, not so much for users")
	flag.StringVarP(&flags.State, "state", "s", "active", "comma separated list of states to filter by. Defaults to 'active' which is a gop shorthand for 'started,finished,delivered,rejected'")
	flag.StringVarP(&flags.User, "user", "u", "me", "Filter by user. 'all' and 'me' are special. Otherwise just a username or initials")
	flag.BoolVarP(&flags.Version, "version", "", false, "Display version information")

	flag.Parse()

	if flags.Version {
		fmt.Printf("Version: %v\n", Version)
	} else if flags.ShellInit {
		fmt.Println(completionScript)
	} else {
		args := flag.Args()
		if len(args) == 0 {
			args = append(args, "usage")
		}

		switch args[0] {
		case "login":
			CommandLogin(flags)
		case "logout":
			CommandLogout(flags)
		case "project":
			if flags.Help {
				HelpProject()
			} else {
				if len(args) == 2 {
					CommandProject(flags, args[1])
				} else {
					CommandProject(flags, "")
				}
			}
		case "config":
			if flags.Help {
				HelpConfig()
			} else {
				CommandConfig(flags)
			}
		case "ls":
			if flags.Help {
				HelpLs()
			} else {
				CommandLs(flags)
			}
		case "current":
			CommandCurrent(flags)
		case "backlog":
			CommandBacklog(flags)
		case "complete":
			if len(args) != 2 {
				fmt.Println("this command takes exactly 1 argument")
			} else {
				CommandComplete(flags, args[1])
			}
		default:
			CommandHelp()
		}
	}
}

func CommandHelp() {
	fmt.Println("Usage: gop [--version] [--help] [-u <user>] [-s <state-list>]")
	fmt.Println("           [-c] [-d] <command> [<args>]")
	fmt.Println("\nAvailable commands are:")
	fmt.Println("   backlog        show all stories in the current projects backlog")
	fmt.Println("   config         set various config options. See config --help for examples")
	fmt.Println("   current        show all stories in the current projects current iteration")
	fmt.Println("   login          provide your pivotal credentials for api access")
	fmt.Println("   ls             list stories in current project")
	fmt.Println("   project [name] show list of projects or set the 'current project' to the one specified")
	fmt.Println("\nNote:")
	fmt.Println("   If you add: eval \"$(gop --shell-init)\"")
	fmt.Println("   to your .zshrc/.bashrc file, you'll get auto-completion for story names")
}

func CommandLogin(flags Flags) {
	var username, password string
	fmt.Printf("Username: ")
	fmt.Scan(&username)
	password, _ = gopass.GetPass("Password: ")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.pivotaltracker.com/services/v5/me", nil)
	req.SetBasicAuth(username, password)
	resp, _ := client.Do(req)
	var currentUser User
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &currentUser)

	Config.CurrentUser = currentUser
	SaveConfig()
}

func CommandLogout(flags Flags) {
	Config.CurrentUser = User{}
	SaveConfig()
}

func HelpConfig() {
	fmt.Println("gop config foo=bar")
	fmt.Println("  Used to set various configuration options")
	fmt.Println("  Currently setable config options include:")
	fmt.Println("    TabCompleteWordCutoff - Number of words to truncate story names to for")
	fmt.Println("                            commands that print/take story names")
}

func CommandConfig(flags Flags) {
	configUpdated := false
	for i := 1; i < len(flag.Args()); i++ {
		x := strings.Split(flag.Args()[i], "=")
		key := x[0]
		value := x[1]
		switch key {
		case "TabCompleteWordCutoff":
			temp, _ := strconv.ParseInt(value, 0, 64)
			Config.TabCompleteWordCutoff = int(temp)
			configUpdated = true
		}
	}
	if configUpdated {
		SaveConfig()
	}
}

func HelpProject() {
	fmt.Println("gop project [project name] - command for working with projects\n")
	fmt.Println("  If not project name is specified, lists all projects that")
	fmt.Println("  the current user belongs to.")
	fmt.Println("  If project name is given, sets the \"current project\" to")
	fmt.Println("  the project specified.")
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

func HelpLs() {
	fmt.Println("gop ls - Lists stories. The default flags are --user me --state active\n")
	fmt.Println("  --user,  -u    # Filter by owner of story. \"me\" refers to the currently")
	fmt.Println("                   logged in user. Can be specified by pivtal account name,")
	fmt.Println("                   email address, or initials")
	fmt.Println("  --state, -s    # Comma-separated list of states from the list unstarted,")
	fmt.Println("                   started, finished, delivered, accepted, rejected.")
	fmt.Println("                   \"active\" is a convinience value that represents")
	fmt.Println("                   \"started,finished,delivered,rejected\"")
	fmt.Println("  --concise, -c  # Lists stories in a more concise format. Useful for")
	fmt.Println("                   piping the results into grep, awk, or other tools")
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

		body, _ := request(fmt.Sprintf("/projects/%v/stories?%s", Config.CurrentProjectId, qParams.Encode()))
		json.Unmarshal(body, &stories)
		printStories(flags, stories)
		recordStories(stories)
	}
}

func CommandComplete(flags Flags, str string) {
	compFilePath := filepath.Join(DbDir, "completions.json")

	completions := make([]Completion, 0)

	compFile, _ := os.OpenFile(compFilePath, os.O_RDWR|os.O_CREATE, 0700)
	defer compFile.Close()

	json.NewDecoder(compFile).Decode(&completions)

	for i := 0; i < len(completions); i++ {
		comp := completions[i]
		if strings.HasPrefix(comp.Text, str) {
			fmt.Println(StoryLongToShort(comp.Text))
		}
	}
}

func StoryShortToLong(shortStory string) string {
	if Config.TabCompleteWordCutoff == 0 {
		return shortStory
	} else {
		compFilePath := filepath.Join(DbDir, "completions.json")

		completions := make([]Completion, 0)

		compFile, _ := os.OpenFile(compFilePath, os.O_RDWR|os.O_CREATE, 0700)
		defer compFile.Close()

		json.NewDecoder(compFile).Decode(&completions)

		for i := 0; i < len(completions); i++ {
			comp := completions[i]
			re := regexp.MustCompile("\\.\\.\\.$")
			if strings.HasPrefix(comp.Text, re.ReplaceAllString(shortStory, "")) {
				return comp.Text
			}
		}
	}
	return ""
}

func StoryLongToShort(longStory string) string {
	var text string
	if Config.TabCompleteWordCutoff == 0 {
		text = longStory
	} else {
		words := strings.Split(longStory, " ")
		text = strings.Join(words[:Config.TabCompleteWordCutoff], " ") + "..."
	}
	return text
}

func CommandCurrent(flags Flags) {
	storiesByIterationScope(flags, "current")
}

func CommandBacklog(flags Flags) {
	storiesByIterationScope(flags, "backlog")
}

func storiesByIterationScope(flags Flags, scope string) {
	iters := make([]Iteration, 0)
	reqStr := fmt.Sprintf("/projects/%v/iterations?scope=%v", Config.CurrentProjectId, scope)
	body, _ := request(reqStr)
	json.Unmarshal(body, &iters)
	if len(iters) != 1 {
		fmt.Println("Got back weird number of iterations")
	} else {
		recordStories(iters[0].Stories)
		printStories(flags, iters[0].Stories)
	}
}

func printStories(flags Flags, stories []Story) {
	if flags.Concise {
		for i := 0; i < len(stories); i++ {
			char := strings.ToUpper(string(stories[i].CurrentState[0]))
			fmt.Printf("%v %v: %v\n", char, stories[i].Id, stories[i].Name)
		}
	} else {
		states := []string{"started", "finished", "delivered", "rejected", "unstarted"}
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

func recordStories(stories []Story) {
	compFilePath := filepath.Join(DbDir, "completions.json")

	completions := make([]Completion, 0)

	compFile, _ := os.OpenFile(compFilePath, os.O_RDWR|os.O_CREATE, 0700)
	defer compFile.Close()

	json.NewDecoder(compFile).Decode(&completions)

	var foundCompletion bool
	for i := 0; i < len(stories); i++ {
		foundCompletion = false
		story := stories[i]
		for j := 0; j < len(completions); j++ {
			comp := completions[j]
			if story.Name == comp.Text {
				comp.LastTouched = time.Now()
				comp.CurrentState = story.CurrentState
				foundCompletion = true
			}
		}
		if !foundCompletion {
			newComp := new(Completion)
			newComp.Id = story.Id
			newComp.Text = story.Name
			newComp.CurrentState = story.CurrentState
			newComp.LastTouched = time.Now()
			completions = append(completions, *newComp)
		}
	}

	newComps := make([]Completion, 0)
	for i := 0; i < len(completions); i++ {
		comp := completions[i]
		if stringInSlice(comp.CurrentState, []string{"started", "finished", "delivered", "rejected"}) || (time.Since(comp.LastTouched)).Hours() < 24*14 {
			newComps = append(newComps, comp)
		}
	}

	fileJson, err := json.Marshal(newComps)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(compFilePath, fileJson, 0622); err != nil {
		panic(err)
	}
}

func stringInSlice(str string, strs []string) bool {
	for i := 0; i < len(strs); i++ {
		if str == strs[i] {
			return true
		}
	}
	return false
}

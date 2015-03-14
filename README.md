# Gop

_a git like client for interacting with pivotal_
===========
Mostly just in the experimentation phase at this point. But who knows, maybe it'll turn into something.

## Usage

Commands

- gop login

  Asks for your username and password and will save your user information including your ApiToken for future requests.

- gop project [project name]

  Lists projects available to the currently logged in user.
  If project name is specified, will set the "current project" to the project with the specified name.

- gop ls

  Lists stories in the current project. By default it lists all stories belonging to the current user

  Add -u or --user to filter by a particular user (you can use their pivotal account, email address, or even just initials).

  Add -s or --state to filter stories to a particular state (comma separated list of states). "all" and "active" are magic words that represent any state and "started, delivered, finished, rejected" respectively. e.g.
  ```bash
  $ gop ls -u EF -s started,rejected
  ```
  Will return all of Elmer Fudd's stories that are currently in a started or rejected state (assuming you work at warner bros.)

There are other commands as well which you can read about in more detail by running gop --help.

## Notes

If after installing gop you add: `eval "$(gop --shell-init)"` to your .zshrc or .bashrc you can get tab-completion for story names.

Most commands can be given --help for help text more specific to that particular command.

Many commands respond to the -c or --concise flag which tends to give a more command-line tool friendly format. For example:
```bash
$ gop ls -u EF -s started,rejected -c
R 1234 Get that wascally wabbit
S 1111 Conduct A Corny Concerto
```
The format of which is "[1-char denoting state of story] [story id] [story name]"

This is still very much an early work in progress. Any and all feedback/advice/opinions are welcome.

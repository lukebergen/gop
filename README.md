# GoPivot

_a git like client for interacting with pivotal_
===========
Mostly just in the experimentation phase at this point. But who knows, maybe it'll turn into something.

## Usage

Commands

- gopivot login

  Asks for your username and password and will save your user information including your ApiToken for future requests.
  This command is run the first time you run gopivot and can be run again if you want to log in with a different user.

- gopivot project [project name]

  Lists projects available to the currently logged in user.
  If project name is specified, will set the "current project" to the project with the specified name.

- gopivot ls

  Lists stories in the current project. By default it lists all stories belonging to the current user

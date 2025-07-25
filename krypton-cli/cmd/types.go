package cmd

import "errors"

type Command interface {
	// execute command. parallel/serial is handled
	Execute() error
	// executes with args via stdin. parallel is ignored now.
	ExecuteWithArgs(interface{}) error
	// parse all args. note command will only get args after command name
	Parse([]string) (Command, error)
	// use for help. prints all args and their help
	PrintDefaults()
	// return interface for stdin input parse.
	// input will be json and each command can return its expected struct
	GetInput() interface{}
}

// convinient data struct for holding commands in a module
type Commands map[string]Command

// module like es, fs etc which is just a bag of commands.
type Module struct {
	Commands
}

// standard errors from command parsing
var (
	ErrMissingArgs = errors.New("Missing arguments")
	ErrParseStdin  = errors.New("Parse standard input")
)

// command run function
type fnRun func()

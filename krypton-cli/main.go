package main

import (
	"errors"
	"fmt"
	"os"

	"cli/auth"
	"cli/cmd"
	"cli/config"
	"cli/dsts"
	"cli/es"
	"cli/fs"
	"cli/iot"
	"cli/logging"
	"cli/ss"
	"cli/util"
)

type module struct {
	name        string
	description string
	cmds        cmd.Commands
}

var (
	ErrConfigLoad     = errors.New("Failed to load config")
	ErrNoModule       = errors.New("No module specified")
	ErrUnknownModule  = errors.New("No such module")
	ErrUnknownCommand = errors.New("No such command")
	modules           = []module{
		module{"auth", "Acquire user tokens", auth.GetCommands()},
		module{"dsts", "Device STS Service", dsts.GetCommands()},
		module{"es", "Enroll Service", es.GetCommands()},
		module{"fs", "File Service", fs.GetCommands()},
		module{"iot", "interact with mqtt server", iot.GetCommands()},
		module{"ss", "Scheduler Service", ss.GetCommands()},
		module{"util", "Utility commands", util.GetCommands()},
	}
	log *logging.Log
)

// parses command args
// determines command to execute
// if all args are present or can be defaulted, executes command
// if args are missing, checks if -stdin is provided, if yes, reads from stdin
func main() {
	log = logging.GetLogger()
	if config.GetSettings() == nil {
		log.Fatal(ErrConfigLoad)
	}
	c, err := parseCommand()
	if errors.Is(err, cmd.ErrParseStdin) {
		log.FatalIf(cmd.ExecuteWithStdin(c))
	} else {
		log.FatalIf(err)
		log.FatalIf(c.Execute())
	}
}

// handles module specific command parsing
func parseCommand() (cmd.Command, error) {
	numArgs := len(os.Args)
	if numArgs < 2 {
		listModules()
		os.Exit(0)
	}
	module := os.Args[1]
	if _, err := getModule(module); err != nil {
		listModules()
		os.Exit(1)
	}
	if numArgs < 3 {
		listCommands(module)
		os.Exit(0)
	}
	return getCommand(module, os.Args[2], os.Args[3:])
}

func getCommand(moduleName, cmdName string, args []string) (cmd.Command, error) {
	m, err := getModule(moduleName)
	if err != nil {
		return nil, err
	}
	return m.cmds.Parse(cmdName, args)
}

func listCommands(moduleName string) {
	m, err := getModule(moduleName)
	if err == nil {
		m.cmds.ListAll()
	}
}

func listModules() {
	fmt.Println("Please specify a module. Supported modules are:")
	for _, m := range modules {
		fmt.Printf("\x1b[0;32m- %v\t(%v) \x1b[0m\n", m.name, m.description)
	}
}

// lookup a module by name
func getModule(name string) (*module, error) {
	for _, m := range modules {
		if m.name == name {
			return &m, nil
		}
	}
	return nil, ErrUnknownModule
}

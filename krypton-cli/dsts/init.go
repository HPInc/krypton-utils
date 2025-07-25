package dsts

import (
	"cli/cmd"
	"cli/logging"
)

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func GetCommands() cmd.Commands {
	return commands
}

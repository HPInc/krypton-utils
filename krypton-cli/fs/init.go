package fs

import (
	"flag"

	"cli/cmd"
	"cli/config"
	"cli/logging"
)

type FileServerFlags struct {
	Server   *string
	JwtToken *string
}

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func (c *FileServerFlags) initServerFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	c.Server = fs.String("server", config.GetSettings().GetAddress("fs", "fs"), "server url")
	c.JwtToken = fs.String("jwt_token", s.GetDeviceToken(), "provide a device token string")
}

func GetCommands() cmd.Commands {
	return commands
}

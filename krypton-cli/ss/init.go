package ss

import (
	"cli/cmd"
	"cli/config"
	"cli/logging"
	"flag"
)

type SchedulerFlags struct {
	HttpServer   *string
	HttpBasePath *string
	DeviceId     *string
	JwtToken     *string
}

type SchedulerBase struct {
	SchedulerFlags
}

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func (b *SchedulerBase) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	b.HttpServer = fs.String("http_server", s.GetAddress("ss", "api"), "http server")
	b.HttpBasePath = fs.String("http_base_path", "api/v1", "http api base path")
	b.DeviceId = fs.String("device_id", s.GetDeviceId(), "device id to be used as client id")
	b.JwtToken = fs.String("jwt_token", s.GetAppToken(), "provide an app token string")
}

func GetCommands() cmd.Commands {
	return commands
}

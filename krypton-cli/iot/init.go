package iot

import (
	"cli/cmd"
	"cli/config"
	"cli/iot/iot_cli"
	"cli/logging"
	"flag"
)

type LoginInput struct {
	DeviceTokenIn string `json:"device_token"`
	DeviceIdIn    string `json:"device_id"`
}

type SchedulerFlags struct {
	IotServer   *string
	DeviceToken *string
	DeviceId    *string
	CertFile    *string
	Timeout     *uint
	Protocol    *string
}

type SchedulerBase struct {
	Client iot_cli.SchedulerClient
	SchedulerFlags
}

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func (b *SchedulerBase) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	b.IotServer = fs.String("iot_server", s.GetAddress("ss", "iot"), "iot endpoint")
	b.DeviceToken = fs.String("device_token", s.GetDeviceToken(), "provide a dsts token string")
	b.DeviceId = fs.String("device_id", s.GetDeviceId(), "device id to be used as client id")
	b.CertFile = fs.String("root_cert_file", s.GetIotRootCertPath(), "root certificate file")
	b.Timeout = fs.Uint("timeout", 5, "timeout in seconds. 0 to wait indefinitely")
	b.Protocol = fs.String("protocol", "mqtt", "protocol: mqtt or ws or wss")
}

func GetCommands() cmd.Commands {
	return commands
}

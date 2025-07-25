package es

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/config"

	"github.com/google/uuid"
)

const (
	CMD_RENEW_ENROLL = "renew_enroll"
)

type RenewEnrollFlags struct {
	deviceId    *common.Uuid
	deviceToken *string
}

type CmdRenewEnroll struct {
	cmd.CmdBase
	EnrollBase
	RenewEnrollInput
}

// to parse and hold std input
type RenewEnrollInput struct {
	DeviceId    uuid.UUID `json:"device_id"`
	DeviceToken string    `json:"device_token"`
}

var (
	deviceToken      = ""
	renewEnrollFlags = RenewEnrollFlags{
		deviceId:    common.NewUUID(),
		deviceToken: &deviceToken,
	}
)

func init() {
	commands[CMD_RENEW_ENROLL] = NewCmdRenewEnroll()
}

func NewCmdRenewEnroll() *CmdRenewEnroll {
	c := CmdRenewEnroll{
		cmd.CmdBase{Name: CMD_RENEW_ENROLL},
		EnrollBase{},
		RenewEnrollInput{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&renewEnrollFlags).initFlags(fs)
	return &c
}

func (u *RenewEnrollFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	fs.Var(u.deviceId, "device_id", "device id (uuid)")
	fs.StringVar(u.deviceToken, "device_token", s.GetDeviceToken(),
		"device token from an enroll")
}

func (u *RenewEnrollFlags) verify() bool {
	if !u.deviceId.IsSet() {
		_ = u.deviceId.Set(config.GetSettings().GetDeviceId())
	}
	return u.deviceId.IsSet() && *u.deviceToken != ""
}

func (c *CmdRenewEnroll) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !(&renewEnrollFlags).verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -device_id and -device_token or specify -stdin for standard input")
			return nil, cmd.ErrMissingArgs
		} else {
			// allow base logic to parse stdin
			err = cmd.ErrParseStdin
		}
	}
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.DeviceId = renewEnrollFlags.deviceId.UUID
	c.DeviceToken = *renewEnrollFlags.deviceToken
	c.RunFunc = c.renewEnroll
	return c, err
}

// main command function
func (c *CmdRenewEnroll) renewEnroll() {
	client := c.Client
	client.SetTokenType("device")
	cr, err := client.RenewEnroll(
		c.DeviceId,
		c.DeviceToken)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(cr))
}

func (c *CmdRenewEnroll) GetInput() interface{} {
	return &c.RenewEnrollInput
}

func (c *CmdRenewEnroll) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

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
	CMD_UNENROLL = "unenroll"
)

type UnenrollFlags struct {
	deviceId    *common.Uuid
	deviceToken *string
}

type CmdUnenroll struct {
	cmd.CmdBase
	EnrollBase
	UnenrollInput
}

var (
	unenrollDeviceToken = ""
	unenrollFlags       = UnenrollFlags{
		deviceId:    common.NewUUID(),
		deviceToken: &unenrollDeviceToken,
	}
)

// to parse and hold std input
type UnenrollInput struct {
	DeviceId    uuid.UUID `json:"device_id"`
	DeviceToken string    `json:"device_token"`
}

func init() {
	commands[CMD_UNENROLL] = NewCmdUnenroll()
}

func NewCmdUnenroll() *CmdUnenroll {
	c := CmdUnenroll{
		cmd.CmdBase{Name: CMD_UNENROLL},
		EnrollBase{},
		UnenrollInput{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&unenrollFlags).initFlags(fs)
	return &c
}

func (u *UnenrollFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	fs.Var(u.deviceId, "device_id", "device id (uuid)")
	fs.StringVar(u.deviceToken, "device_token", s.GetDeviceToken(),
		"device token from an enroll")
}

func (u *UnenrollFlags) verify() bool {
	if !u.deviceId.IsSet() {
		_ = u.deviceId.Set(config.GetSettings().GetDeviceId())
	}
	return u.deviceId.IsSet() && *u.deviceToken != ""
}

func (c *CmdUnenroll) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)

	if !(&unenrollFlags).verify() {
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
	c.DeviceId = unenrollFlags.deviceId.UUID
	c.DeviceToken = *unenrollFlags.deviceToken
	c.RunFunc = c.unenroll
	return c, err
}

// main worker function
func (c *CmdUnenroll) unenroll() {
	client := c.Client
	cr, err := client.Unenroll(
		c.DeviceId,
		c.DeviceToken)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(cr))
}

func (c *CmdUnenroll) GetInput() interface{} {
	return &c.UnenrollInput
}

func (c *CmdUnenroll) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

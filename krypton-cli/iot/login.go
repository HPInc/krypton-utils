package iot

import (
	"cli/cmd"
	"cli/iot/iot_cli"
)

const (
	CmdLoginName = "login"
)

type CmdLogin struct {
	cmd.CmdBase
	SchedulerBase
	LoginInput
}

func init() {
	commands[CmdLoginName] = NewCmdLogin()
}

func NewCmdLogin() *CmdLogin {
	c := CmdLogin{
		cmd.CmdBase{
			Name: CmdLoginName,
		},
		SchedulerBase{},
		LoginInput{},
	}
	fs := c.BaseInitFlags()
	(&c.SchedulerBase).initFlags(fs)
	return &c
}

func (c *CmdLogin) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if c.Stdin {
		err = cmd.ErrParseStdin
	} else {
		c.DeviceTokenIn = *c.DeviceToken
		c.DeviceIdIn = *c.DeviceId
	}
	c.RunFunc = c.login
	return c, err
}

func (c *CmdLogin) login() {
	cli := iot_cli.NewClient(
		*c.IotServer,
		*c.CertFile,
		c.DeviceIdIn,
		c.DeviceTokenIn,
		*c.Timeout,
		*c.Protocol,
	)
	if err := cli.Login(); err != nil {
		log.Fatal("Login error: ", err)
	}
}

func (c *CmdLogin) GetInput() interface{} {
	return &c.LoginInput
}

func (c *CmdLogin) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

package es

import (
	"fmt"

	"cli/cmd"
	"cli/common"
)

const (
	cmdDeleteEnrollToken = "delete_enroll_token"
)

type CmdDeleteEnrollToken struct {
	cmd.CmdBase
	EnrollBase
}

func init() {
	commands[cmdDeleteEnrollToken] = NewCmdDeleteEnrollToken()
}

func NewCmdDeleteEnrollToken() *CmdDeleteEnrollToken {
	c := CmdDeleteEnrollToken{
		cmd.CmdBase{
			Name: cmdDeleteEnrollToken,
		},
		EnrollBase{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	return &c
}

func (c *CmdDeleteEnrollToken) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.RunFunc = c.deleteEnrollToken
	return c, nil
}

func (c *CmdDeleteEnrollToken) deleteEnrollToken() {
	client := c.Client
	er, err := client.DeleteEnrollToken()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(er))
}

func (c *CmdDeleteEnrollToken) GetInput() interface{} {
	return nil
}

func (c *CmdDeleteEnrollToken) ExecuteWithArgs(i interface{}) error {
	return nil
}

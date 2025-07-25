package es

import (
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/config"
)

const (
	cmdCreateEnrollToken = "create_enroll_token"
)

type CmdCreateEnrollToken struct {
	cmd.CmdBase
	EnrollBase
}

func init() {
	commands[cmdCreateEnrollToken] = NewCmdCreateEnrollToken()
}

func NewCmdCreateEnrollToken() *CmdCreateEnrollToken {
	c := CmdCreateEnrollToken{
		cmd.CmdBase{
			Name: cmdCreateEnrollToken,
		},
		EnrollBase{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	return &c
}

func (c *CmdCreateEnrollToken) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.RunFunc = c.createEnrollToken
	return c, nil
}

func (c *CmdCreateEnrollToken) createEnrollToken() {
	client := c.Client
	er, err := client.CreateEnrollToken()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	result := common.GetJsonString(er)
	// update cli cache with enroll_token
	// only write cache on single count iterations
	if c.Count == 1 {
		if err = config.GetBulkEnrollCache().Update([]byte(result)); err != nil {
			log.Debug("failed to cache bulk enroll token")
		}
	}
	fmt.Println(result)
}

func (c *CmdCreateEnrollToken) GetInput() interface{} {
	return nil
}

func (c *CmdCreateEnrollToken) ExecuteWithArgs(i interface{}) error {
	return nil
}

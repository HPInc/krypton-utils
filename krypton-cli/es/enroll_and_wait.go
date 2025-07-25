package es

import (
	"fmt"

	"cli/cmd"
	"cli/common"
)

const (
	cmdEnrollAndWait = "enroll_and_wait"
)

type CmdEnrollAndWait struct {
	cmd.CmdBase
	EnrollBase
}

func init() {
	commands[cmdEnrollAndWait] = NewCmdEnrollAndWait()
}

func NewCmdEnrollAndWait() *CmdEnrollAndWait {
	c := CmdEnrollAndWait{
		cmd.CmdBase{Name: cmdEnrollAndWait},
		EnrollBase{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&enrollFlags).initFlags(fs)
	return &c
}

func (c *CmdEnrollAndWait) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.Client.HardwareHash = *enrollFlags.hardwareHash
	c.Client.BulkEnrollToken = *enrollFlags.bulkEnrollToken
	c.RunFunc = c.enrollAndWait
	return c, nil
}

func (c CmdEnrollAndWait) enrollAndWait() {
	client := c.Client
	cr, err := client.GetDeviceCertificate()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(cr))
}

func (c CmdEnrollAndWait) GetInput() interface{} {
	return nil
}

func (c CmdEnrollAndWait) ExecuteWithArgs(i interface{}) error {
	return nil
}

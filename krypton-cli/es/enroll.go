package es

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/config"
)

const (
	cmdEnroll = "enroll"
)

type EnrollFlags struct {
	bulkEnrollToken *string
	hardwareHash    *string
}

type CmdEnroll struct {
	cmd.CmdBase
	EnrollBase
}

var (
	hardwareHash    = ""
	bulkEnrollToken = ""
	enrollFlags     = EnrollFlags{
		hardwareHash:    &hardwareHash,
		bulkEnrollToken: &bulkEnrollToken,
	}
)

func init() {
	commands[cmdEnroll] = NewCmdEnroll()
}

func NewCmdEnroll() *CmdEnroll {
	c := CmdEnroll{
		cmd.CmdBase{
			Name: cmdEnroll,
		},
		EnrollBase{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&enrollFlags).initFlags(fs)
	return &c
}

func (e *EnrollFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	fs.StringVar(e.hardwareHash, "hardware_hash", "", "device hardware hash")
	fs.StringVar(e.bulkEnrollToken, "bulk_enroll_token",
		s.GetBulkEnrollToken(), "bulk enroll token")
}

func (c *CmdEnroll) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.Client.HardwareHash = *enrollFlags.hardwareHash
	c.Client.BulkEnrollToken = *enrollFlags.bulkEnrollToken
	c.RunFunc = c.enroll
	return c, nil
}

func (c *CmdEnroll) enroll() {
	client := c.Client
	cr, err := client.Enroll()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(cr))
}

func (c *CmdEnroll) GetInput() interface{} {
	return nil
}

func (c *CmdEnroll) ExecuteWithArgs(i interface{}) error {
	return nil
}

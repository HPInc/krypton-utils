package es

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/config"
)

const (
	CmdGetStatusName = "get_status"
)

type EnrollStatusFlags struct {
	enrollId        *common.Uuid
	bulkEnrollToken *string
}

// to parse and hold std input
type EnrollStatusInput struct {
	EnrollId string `json:"id"`
}

type CmdGetStatus struct {
	cmd.CmdBase
	EnrollBase
	EnrollStatusInput
}

var (
	enrollStatusFlags = EnrollStatusFlags{
		enrollId:        common.NewUUID(),
		bulkEnrollToken: &bulkEnrollToken,
	}
)

func init() {
	commands[CmdGetStatusName] = NewCmdGetStatus()
}

func NewCmdGetStatus() *CmdGetStatus {
	c := CmdGetStatus{
		cmd.CmdBase{Name: CmdGetStatusName},
		EnrollBase{},
		EnrollStatusInput{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&enrollStatusFlags).initFlags(fs)
	return &c
}

func (f *EnrollStatusFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	fs.Var(f.enrollId, "enroll_id", "enroll id (uuid)")
	fs.StringVar(f.bulkEnrollToken, "bulk_enroll_token",
		s.GetBulkEnrollToken(), "bulk enroll token")
}

func (f *EnrollStatusFlags) verify() bool {
	return f.enrollId.IsSet()
}

func (c *CmdGetStatus) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)

	if !(&enrollStatusFlags).verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -enroll_id or specify -stdin for standard input")
			return nil, cmd.ErrMissingArgs
		} else {
			// allow base logic to parse stdin
			err = cmd.ErrParseStdin
		}
	}
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.EnrollId = enrollStatusFlags.enrollId.String()
	c.RunFunc = c.getStatus
	return c, err
}

// main worker function
func (c *CmdGetStatus) getStatus() {
	client := c.Client
	client.BulkEnrollToken = *enrollStatusFlags.bulkEnrollToken
	cr, err := client.GetStatus(c.EnrollId)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(common.GetJsonString(cr))
}

func (c *CmdGetStatus) GetInput() interface{} {
	return &c.EnrollStatusInput
}

func (c *CmdGetStatus) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

package fs

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/fs/restcli"
)

const (
	CMD_GET_FILE_DETAILS = "get_file_details"
)

type GetFileDetailsFlags struct {
	FileId *int
}

type GetFileIdInput struct {
	FileId int `json:"file_id"`
}

type CmdGetFileDetails struct {
	cmd.CmdBase
	FileServerFlags
	GetFileIdInput
}

var (
	getFileDetailsFlags = GetFileDetailsFlags{}
)

func init() {
	commands[CMD_GET_FILE_DETAILS] = NewCmdGetFileDetails()
}

func NewCmdGetFileDetails() *CmdGetFileDetails {
	c := CmdGetFileDetails{
		cmd.CmdBase{Name: CMD_GET_FILE_DETAILS},
		FileServerFlags{},
		GetFileIdInput{},
	}
	fs := c.BaseInitFlags()
	(&c.FileServerFlags).initServerFlags(fs)
	(&getFileDetailsFlags).initFlags(fs)
	return &c
}

func (u *GetFileDetailsFlags) initFlags(fs *flag.FlagSet) {
	u.FileId = fs.Int("file_id", 1, "file id")
}

func (c *CmdGetFileDetails) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if c.Stdin {
		err = cmd.ErrParseStdin
	} else {
		c.FileId = *getFileDetailsFlags.FileId
	}
	c.RunFunc = c.getFileDetails
	return c, err
}

func (c *CmdGetFileDetails) getFileDetails() {
	r, e := restcli.NewFileDetailsClient(*c.Server, *c.JwtToken).Execute(c.FileId)
	if e != nil {
		log.Fatal("Error: ", e)
	} else {
		fmt.Println(common.GetJsonString(r.Data))
	}
}

func (c *CmdGetFileDetails) GetInput() interface{} {
	return &c.GetFileIdInput
}

func (c *CmdGetFileDetails) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

package fs

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/fs/restcli"
)

const (
	CMD_GET_DOWNLOAD_URL = "get_download_url"
)

type DownloadUrlFlags struct {
	FileId *int
}

type CmdGetDownloadUrl struct {
	cmd.CmdBase
	FileServerFlags
	DownloadUrlFlags
}

type GetDownloadUrlResult struct {
	Url string `json:"url"`
}

func init() {
	commands[CMD_GET_DOWNLOAD_URL] = NewCmdGetDownloadUrl()
}

func NewCmdGetDownloadUrl() *CmdGetDownloadUrl {
	c := CmdGetDownloadUrl{
		cmd.CmdBase{Name: CMD_GET_DOWNLOAD_URL},
		FileServerFlags{},
		DownloadUrlFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.FileServerFlags).initServerFlags(fs)
	(&c.DownloadUrlFlags).initFlags(fs)
	return &c
}

func (u *DownloadUrlFlags) initFlags(fs *flag.FlagSet) {
	u.FileId = fs.Int("file_id", 1, "file id")
}

func (c *CmdGetDownloadUrl) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	c.RunFunc = c.getDownloadUrl
	return c, nil
}

func (c *CmdGetDownloadUrl) getDownloadUrl() {
	r, e := restcli.NewDownloadUrlClient(*c.Server, *c.JwtToken).Execute(
		*c.DownloadUrlFlags.FileId)
	if e != nil {
		log.Fatal("Error: ", e)
	} else {
		fmt.Println(common.GetJsonString(r))
	}
}

func (c *CmdGetDownloadUrl) GetInput() interface{} {
	return nil
}

func (c *CmdGetDownloadUrl) ExecuteWithArgs(i interface{}) error {
	return nil
}

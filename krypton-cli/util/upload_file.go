package util

import (
	"cli/cmd"
	"cli/common"
	"flag"
)

const (
	cmdUploadFileName = "upload_file"
)

// additional flags for upload_file
type UploadFileFlags struct {
	url  *string
	file *string
}

// to parse and hold std input
type UploadFileInput struct {
	Url  string `json:"url"`
	File string `json:"file"`
}

type CmdUploadFile struct {
	cmd.CmdBase
	UploadFileInput
}

var uploadFileFlags = UploadFileFlags{}

func init() {
	commands[cmdUploadFileName] = NewCmdUploadFile()
}

func NewCmdUploadFile() *CmdUploadFile {
	c := &CmdUploadFile{
		cmd.CmdBase{
			Name: cmdUploadFileName,
		},
		UploadFileInput{},
	}
	fs := c.BaseInitFlags()
	(&uploadFileFlags).initFlags(fs)
	return c
}

func (f *UploadFileFlags) initFlags(fs *flag.FlagSet) {
	f.url = fs.String("url", "localhost", "upload url")
	f.file = fs.String("file", "", "file to upload")
}

func (c *CmdUploadFile) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !uploadFileFlags.verify() {
		if !c.Stdin {
			log.Error("Please provide -server and -file")
			return nil, cmd.ErrMissingArgs
		}
	} else {
		c.Url = *uploadFileFlags.url
		c.File = *uploadFileFlags.file
	}
	// allow stdin to override any cached input
	if c.Stdin {
		err = cmd.ErrParseStdin
	}
	c.RunFunc = c.uploadFile
	return c, err
}

func (f *UploadFileFlags) verify() bool {
	return *f.url != "" && *f.file != ""
}

func (c *CmdUploadFile) GetInput() interface{} {
	return &c.UploadFileInput
}

func (c *CmdUploadFile) uploadFile() {
	if err := common.UploadFile(c.File, nil, c.Url); err != nil {
		log.Fatal("Error: ", err)
	}
}

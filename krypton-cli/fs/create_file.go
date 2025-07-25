package fs

import (
	"flag"
	"fmt"
	"os"

	"cli/cmd"
	"cli/common"
	"cli/config"
	"cli/fs/restcli"
)

const (
	CMD_CREATE_FILE = "create_file"
)

type CreateFileFlags struct {
	FileName    *string
	MaxFileSize *int
	Encrypted   *bool
}

type CmdCreateFile struct {
	cmd.CmdBase
	FileServerFlags
	CreateFileFlags
}

func init() {
	commands[CMD_CREATE_FILE] = NewCmdCreateFile()
}

func NewCmdCreateFile() *CmdCreateFile {
	c := CmdCreateFile{
		cmd.CmdBase{Name: CMD_CREATE_FILE},
		FileServerFlags{},
		CreateFileFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.FileServerFlags).initServerFlags(fs)
	(&c.CreateFileFlags).initFlags(fs)
	return &c
}

func (u *CreateFileFlags) initFlags(fs *flag.FlagSet) {
	u.MaxFileSize = fs.Int("max_file_size", 100, "maximum size of file upload")
	u.FileName = fs.String("filename", "", "file name with path")
	u.Encrypted = fs.Bool("encrypted", false, "pick from encrypted data in config")
}

func (c *CmdCreateFile) verify() bool {
	return *c.JwtToken != ""
}

func (c *CmdCreateFile) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	if !c.verify() {
		return nil, cmd.ErrMissingArgs
	}
	c.RunFunc = c.createFile
	return c, nil
}

// do not remove files that are not auto created
func (c *CmdCreateFile) removeFile(fileName string) {
	// if encrypted, data is fetched from config
	// and no disk files are involved
	if *c.FileName != "" || *c.Encrypted {
		return
	}
	if err := os.Remove(fileName); err != nil {
		log.Error("Could not remove file: ", fileName, "Error: ", err)
	}
}

func (c *CmdCreateFile) createFile() {
	var fileData []byte
	uploadFileName := *c.FileName
	if uploadFileName == "" {
		uploadFileName = createNewFile(*c.CreateFileFlags.MaxFileSize)
	} else if *c.Encrypted {
		fileData = decryptContents(uploadFileName)
	}
	p := restcli.UploadUrlParams{
		FileName: uploadFileName,
		FileData: fileData,
	}
	defer c.removeFile(p.FileName)
	r, err := restcli.NewUploadUrlClient(*c.Server, *c.JwtToken).Execute(p)
	if err != nil {
		log.Fatal("Error: ", err)
	} else {
		url, err := r.GetUploadUrl()
		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			if err = common.UploadFile(p.FileName, p.FileData, url); err != nil {
				log.Fatal("Error: ", err)
			}
			fmt.Println(common.GetJsonString(r.Data))
		}
	}
}

func decryptContents(name string) []byte {
	data := config.GetEncryptedTestData(name)
	if data == nil {
		log.Fatalf("Error: no such data: %s. check config.", name)
	}
	bytes, err := common.Decrypt(data.Data, data.Key)
	if err != nil {
		log.Fatalf("Error: %v. Could not decrypt: %s.", err, name)
	}
	return bytes
}

func createNewFile(maxSize int) string {
	prefix := "fs_cli_upload_*"
	filename, _ := common.TempFileWithRandomContents(prefix, maxSize)
	return filename
}

func (c *CmdCreateFile) GetInput() interface{} {
	return nil
}

func (c *CmdCreateFile) ExecuteWithArgs(i interface{}) error {
	return nil
}

package fs

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/fs/restcli"
)

const (
	CMD_GET_UPLOAD_URL = "get_upload_url"
)

type UploadUrlFlags struct {
	TenantId *common.Uuid
	DeviceId *common.Uuid
	FileName *string
}

type CmdGetUploadUrl struct {
	cmd.CmdBase
	FileServerFlags
	UploadUrlFlags
}

func init() {
	commands[CMD_GET_UPLOAD_URL] = NewCmdGetUploadUrl()
}

func NewCmdGetUploadUrl() *CmdGetUploadUrl {
	c := CmdGetUploadUrl{
		cmd.CmdBase{Name: CMD_GET_UPLOAD_URL},
		FileServerFlags{},
		UploadUrlFlags{
			TenantId: common.NewUUID(),
			DeviceId: common.NewUUID(),
		},
	}
	fs := c.BaseInitFlags()
	(&c.FileServerFlags).initServerFlags(fs)
	(&c.UploadUrlFlags).initFlags(fs)
	return &c
}

func (u *UploadUrlFlags) initFlags(fs *flag.FlagSet) {
	fs.Var(u.TenantId, "tenant_id", "tenant id (uuid)")
	fs.Var(u.DeviceId, "device_id", "device id (uuid)")
	u.FileName = fs.String("filename", "/tmp/1", "file name with path")
}

func (u *UploadUrlFlags) verify() bool {
	if !u.TenantId.IsSet() {
		u.TenantId.SetDefault()
	}
	if !u.DeviceId.IsSet() {
		u.DeviceId.SetDefault()
	}
	return u.TenantId.IsSet() && u.DeviceId.IsSet()
}

func (c *CmdGetUploadUrl) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	if !(&c.UploadUrlFlags).verify() {
		return c, fmt.Errorf("please provide tenant_id and device_id")
	}
	c.RunFunc = c.getUploadUrl
	return c, nil
}

// exits on error
func (c *CmdGetUploadUrl) getUploadUrl() {
	p := restcli.UploadUrlParams{
		FileName: *c.UploadUrlFlags.FileName,
	}
	r, e := restcli.NewUploadUrlClient(*c.Server, *c.JwtToken).Execute(p)
	if e != nil {
		log.Fatal("Error: ", e)
	} else {
		fmt.Println(common.GetJsonString(r.Data))
	}
}

func (c *CmdGetUploadUrl) GetInput() interface{} {
	return nil
}

func (c *CmdGetUploadUrl) ExecuteWithArgs(i interface{}) error {
	return nil
}

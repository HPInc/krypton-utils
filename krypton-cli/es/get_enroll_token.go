package es

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"cli/cmd"
	"cli/common"
	"cli/config"
	"cli/es/enroll"
)

const (
	cmdGetEnrollToken = "get_enroll_token"
)

type GetEnrollTokenFlags struct {
	TenantId *string
	AppToken *string
}

type GetTenantIdInput struct {
	TenantId string `json:"tenant_id"`
}

type CmdGetEnrollToken struct {
	cmd.CmdBase
	EnrollBase
	GetTenantIdInput
}

var (
	getEnrollTokenFlags = GetEnrollTokenFlags{}
)

func init() {
	commands[cmdGetEnrollToken] = NewCmdGetEnrollToken()
}

func NewCmdGetEnrollToken() *CmdGetEnrollToken {
	c := CmdGetEnrollToken{
		cmd.CmdBase{
			Name: cmdGetEnrollToken,
		},
		EnrollBase{},
		GetTenantIdInput{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&getEnrollTokenFlags).initFlags(fs)
	return &c
}

func (e *GetEnrollTokenFlags) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	e.TenantId = fs.String("tenant_id", "", "tenant id for enroll token")
	e.AppToken = fs.String("app_token", s.GetAppToken(),
		"app token from an auth app_token request")
}

func (c *CmdGetEnrollToken) verify() bool {
	return *getEnrollTokenFlags.TenantId != "" &&
		*getEnrollTokenFlags.AppToken != ""
}

func (c *CmdGetEnrollToken) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !c.verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -tenant_id, -app_token or specify -stdin for standard input")
			return nil, cmd.ErrMissingArgs
		} else {
			err = cmd.ErrParseStdin
		}
	} else {
		c.TenantId = *getEnrollTokenFlags.TenantId
	}
	c.RunFunc = c.getEnrollToken
	return c, err
}

func (c *CmdGetEnrollToken) getEnrollToken() {
	if c.TenantId == "" {
		log.Fatalf("Please specify -tenant_id or use -stdin")
	}
	url := fmt.Sprintf("%s/%s/enroll_token/%s", *c.Flags.Server, c.ApiBasePath, c.TenantId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("failed to get enroll_token: %v\n", err)
	}
	req.Header.Add("X-HP-Token-Type", "app")
	req.Header.Add("Authorization", "Bearer "+*getEnrollTokenFlags.AppToken)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Get enroll token for %s failed: %v\n", c.TenantId, err)
	}
	log.HttpResponse(resp, nil)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Get enroll token failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read enrollment status check response. Error: %v\n", err)
	}
	log.HttpResponse(resp, data)
	er := enroll.EnrollTokenResponse{}
	err = json.Unmarshal(data, &er)
	if err != nil {
		log.Fatalf("Failed to unmarshal enroll token response. Error: %v\n", err)
	}
	fmt.Println(common.GetJsonString(er))
}

func (c *CmdGetEnrollToken) GetInput() interface{} {
	return &c.GetTenantIdInput
}

func (c *CmdGetEnrollToken) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

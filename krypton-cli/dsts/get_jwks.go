package dsts

import (
	"cli/cmd"
	"cli/common"
	"cli/config"
	"flag"
	"fmt"
	"io"
	"net/http"
)

const (
	CmdGetJwksName = "keys"
)

// additional flags for get jwks
type getJwksFlags struct {
	server *string
}

type CmdGetJwks struct {
	cmd.CmdBase
	getJwksFlags
}

func init() {
	commands[CmdGetJwksName] = NewCmdGetJwks()
}

func NewCmdGetJwks() *CmdGetJwks {
	c := &CmdGetJwks{
		cmd.CmdBase{
			Name: CmdGetJwksName,
		},
		getJwksFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.getJwksFlags).initFlags(fs)
	return c
}

func (f *getJwksFlags) initFlags(fs *flag.FlagSet) {
	f.server = fs.String("server",
		config.GetSettings().GetAddress("es", "dsts"),
		"dsts server address")
}

func (c *CmdGetJwks) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !c.verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -server")
			return nil, cmd.ErrMissingArgs
		} else {
			err = cmd.ErrParseStdin
		}
	}
	c.RunFunc = c.getJwks
	return c, err
}

func (c *CmdGetJwks) verify() bool {
	return *c.server != ""
}

func (c *CmdGetJwks) getJwks() {
	req, err := http.NewRequest(http.MethodGet, c.getJwksUrl(), nil)
	if err != nil {
		log.Fatalf("failed to get jwks request: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("get jwks failed: %v\n", err)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read jwks response. Error: %v\n", err)
	}
	log.HttpResponse(resp, result)
	fmt.Println(string(result))
}

func (c *CmdGetJwks) getJwksUrl() string {
	return fmt.Sprintf("%s/%s/keys", *c.server, c.ApiBasePath)
}

package util

import (
	"cli/cmd"
	"cli/common"
	"flag"
	"net/http"
	"strconv"
	"strings"
)

const (
	CmdRetryUrlName = "retry_url"
)

// additional flags for wait for server
type RetryUrlFlags struct {
	url          *string
	status       *int
	ignoreStatus *string
}

type CmdRetryUrl struct {
	cmd.CmdBase
	RetryUrlFlags
}

func init() {
	commands[CmdRetryUrlName] = NewCmdRetryUrl()
}

func NewCmdRetryUrl() *CmdRetryUrl {
	c := &CmdRetryUrl{
		cmd.CmdBase{
			Name: CmdRetryUrlName,
		},
		RetryUrlFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.RetryUrlFlags).initFlags(fs)
	return c
}

func (f *RetryUrlFlags) initFlags(fs *flag.FlagSet) {
	f.url = fs.String("url", "", "url to retry till desired status")
	f.status = fs.Int("status", 200, "desired status for success")
	f.ignoreStatus = fs.String("ignore_status", "", "comma separated status strings to ignore and continue")
}

func (c *CmdRetryUrl) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	if !c.verify() {
		if !c.Stdin {
			log.Error(
				"Please provide -url")
			return nil, cmd.ErrMissingArgs
		} else {
			err = cmd.ErrParseStdin
		}
	}
	c.RunFunc = c.waitForUrl
	return c, err
}

func (c *CmdRetryUrl) verify() bool {
	return *c.url != "" && *c.status != 0
}

// return a list of status codes to ignore
// this is needed because RetriableClient will
// fail on status codes like 404 but if we are
// waiting for a pending url like files/123, ignoring
// 404 is useful.
func (c *CmdRetryUrl) getIgnoreStatusList() []int {
	ignoreStatusList := []int{}
	for _, ignore := range strings.Split(*c.ignoreStatus, ",") {
		i, _ := strconv.Atoi(ignore)
		if i > 0 {
			ignoreStatusList = append(ignoreStatusList, i)
		}
	}
	return ignoreStatusList
}

func (c *CmdRetryUrl) waitForUrl() {
	url := *c.url
	expectedStatus := *c.status
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error waiting for %s: %v\n", url, err)
	}
	common.AddUserAgentHeader(req)
	client := common.RetriableClient(c.RetryCount)
	client.IgnoreStatusList = c.getIgnoreStatusList()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error waiting for %s: %v\n", url, err)
	}
	returnStatus := resp.StatusCode
	if returnStatus != expectedStatus {
		log.Fatalf("Url: %s - Expected status %d. Got %d\n",
			url, expectedStatus, returnStatus)
	}
	log.Printf("Status = %d", returnStatus)
}

package util

import (
	"cli/cmd"
	"cli/common"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	CmdWaitForServerName = "wait_for_server"
)

// additional flags for wait for server
type WaitForServerFlags struct {
	scheme *string
	server *string
	port   *uint
}

type CmdWaitForServer struct {
	cmd.CmdBase
	WaitForServerFlags
}

func init() {
	commands[CmdWaitForServerName] = NewCmdWaitForServer()
}

func NewCmdWaitForServer() *CmdWaitForServer {
	c := &CmdWaitForServer{
		cmd.CmdBase{
			Name: CmdWaitForServerName,
		},
		WaitForServerFlags{},
	}
	fs := c.BaseInitFlags()
	(&c.WaitForServerFlags).initFlags(fs)
	return c
}

func (f *WaitForServerFlags) initFlags(fs *flag.FlagSet) {
	f.scheme = fs.String("scheme", "http", "scheme")
	f.server = fs.String("server", "localhost", "server address")
	f.port = fs.Uint("port", 80, "port")
}

func (c *CmdWaitForServer) Parse(args []string) (cmd.Command, error) {
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
	c.RunFunc = c.waitForServer
	return c, err
}

func (c *CmdWaitForServer) verify() bool {
	return *c.server != "" && *c.port != 0
}

func (c *CmdWaitForServer) waitForServer() {
	if strings.HasPrefix(*c.scheme, "http") {
		c.doHttpRetry()
	} else {
		c.doTcpRetry()
	}
}

func connect(addr *net.TCPAddr) func() bool {
	return func() bool {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}
}

func (c *CmdWaitForServer) doTcpRetry() {
	scheme := *c.scheme
	if scheme != "" {
		scheme = scheme + "://"
	}
	srv := fmt.Sprintf("%s%s:%d", scheme, *c.server, *c.port)
	log.Debugf("Connecting to %s\n", srv)
	addr, err := net.ResolveTCPAddr("tcp", srv)
	if err != nil {
		log.Fatalf("Error resolving %s: %v\n", srv, err)
	}
	common.RetryWait(c.RetryCount, connect(addr))
}

func (c *CmdWaitForServer) doHttpRetry() {
	path := c.ApiBasePath
	if path != "" {
		path = "/" + path
	}
	url := fmt.Sprintf("%s://%s:%d%s", *c.scheme, *c.server, *c.port, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error waiting for %s: %v\n", url, err)
	}
	common.AddUserAgentHeader(req)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error waiting for %s: %v\n", url, err)
	}
	log.Printf("Status = %d", resp.StatusCode)
}

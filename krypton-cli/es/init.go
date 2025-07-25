package es

import (
	"flag"
	"fmt"

	"cli/cmd"
	"cli/config"
	"cli/es/enroll"
	"cli/logging"
)

type EnrollServerFlags struct {
	Server           *string
	DSTSServer       *string
	TokenServer      *string
	TokenType        *string
	JWTToken         *string
	ManagementServer *string
}

type EnrollBase struct {
	Client enroll.EnrollClient
	Flags  EnrollServerFlags
}

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func (b *EnrollBase) initFlags(fs *flag.FlagSet) {
	s := config.GetSettings()
	b.Flags.Server = fs.String("server", s.GetAddress("es", "es"), "enroll server")
	b.Flags.TokenServer = fs.String("token_server", getTokenServer(s), "test token server")
	b.Flags.DSTSServer = fs.String("dsts_server", s.GetAddress("es", "dsts"), "dsts server")
	b.Flags.TokenType = fs.String("token_type", s.GetTokenType(), "token type")
	b.Flags.JWTToken = fs.String("jwt_token", s.GetAccessToken(), "provide a jwt token string")
	b.Flags.ManagementServer = fs.String("mgmt_server", s.GetManagementServer(),
		"management server for this client (hpcem, hpconnect)")
}

func (b *EnrollBase) initClient(retryCount uint, apiBasePath string) {
	b.Client = enroll.EnrollClient{
		EnrollUrl:        fmt.Sprintf("%s/%s", *b.Flags.Server, apiBasePath),
		DSTSServer:       fmt.Sprintf("%s/%s", *b.Flags.DSTSServer, apiBasePath),
		JWTToken:         *b.Flags.JWTToken,
		TokenServer:      *b.Flags.TokenServer,
		ManagementServer: *b.Flags.ManagementServer,
	}
	b.Client.SetTokenType(*b.Flags.TokenType)
	b.Client.SetRetryCount(retryCount)
}

func GetCommands() cmd.Commands {
	return commands
}

func getTokenServer(c *config.Config) string {
	return fmt.Sprintf("%s/%s", c.GetAddress("es", "auth"), "api/v1/token")
}

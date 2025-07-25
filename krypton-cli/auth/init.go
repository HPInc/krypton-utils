package auth

import (
	"flag"

	"cli/cmd"
	"cli/config"
	"cli/logging"
)

type AuthTokenType string

const (
	AuthTokenTypeApp      = "app"
	AuthTokenTypeAuth     = "auth"
	TokenProviderOneCloud = "onecloud"
)

type AuthFlags struct {
	Server     *string
	RetryDelay *uint
	TokenType  *string
}

type AuthBase struct {
	Flags AuthFlags
}

var (
	commands = make(cmd.Commands)
	log      = logging.GetLogger()
)

func (b *AuthBase) initFlags(fs *flag.FlagSet, t AuthTokenType) {
	s := config.GetSettings()
	b.Flags.Server = fs.String("server", s.GetAddress("auth", string(t)), "server")
	b.Flags.RetryDelay = fs.Uint("retry_wait", 5, "number of seconds to wait between retries")
	b.Flags.TokenType = fs.String("token_type", s.CurrentProfile.Auth.TokenType, "token_type")
}

func GetCommands() cmd.Commands {
	return commands
}

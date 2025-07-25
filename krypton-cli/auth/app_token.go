package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"cli/cmd"
	"cli/config"
)

const (
	CmdAppTokenName = "app_token"
)

type AppTokenFlags struct {
	AppId          *string
	PrivateKeyFile *string
}

type CmdAppToken struct {
	cmd.CmdBase
	AuthBase
	AppTokenFlags
	PrivateKey *rsa.PrivateKey
}

func init() {
	commands[CmdAppTokenName] = NewCmdAppToken()
}

func NewCmdAppToken() *CmdAppToken {
	c := CmdAppToken{
		cmd.CmdBase{
			Name: CmdAppTokenName,
		},
		AuthBase{},
		AppTokenFlags{},
		nil,
	}
	fs := c.BaseInitFlags()
	(&c.AuthBase).initFlags(fs, AuthTokenTypeApp)
	(&c.AppTokenFlags).initFlags(fs)
	return &c
}

func (f *AppTokenFlags) initFlags(fs *flag.FlagSet) {
	f.AppId = fs.String("app_id", "de7e595f-9aca-4334-9f47-2352d00acace",
		"app id of a registered app in dsts")
	f.PrivateKeyFile = fs.String("pk_file", "/tmp/privateKey.pem",
		"private key file of application registering for app token")
}

func (c *CmdAppToken) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	c.RunFunc = c.getAppToken
	return c, nil
}

// flow to get an app token
// Similar flow as in device token below.
// replace deviceauth with appauth and device_id with app_id
// see https://rndwiki.inc.hpicorp.net/confluence/display/Hptm01/DSTS+API%3A+Device+authentication
func (c *CmdAppToken) getAppToken() {
	// get challenge code
	challenge, err := c.getChallengeCode()
	if err != nil {
		log.Fatalf("AppToken: Failed to get challenge code. AppId: %s, Error: %v\n",
			*c.AppId, err)
	}
	c.PrivateKey, err = getPrivateKey(c.PrivateKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	assertion, err := c.getAssertion(challenge)
	if err != nil {
		log.Fatalf("Failed to generate signed client assertion. AppId: %s, Error: %v",
			*c.AppId, err)
	}
	log.Debug("assertion: ", assertion)
	token, err := c.getAppTokenFromSTS(assertion)
	if err != nil {
		log.Fatal(err)
	}
	// update cli cache with token
	if err = config.GetAppTokenCache().Update(token); err != nil {
		log.Debug("failed to cache app access token")
	}
	fmt.Println(token)
}

func getPrivateKey(pkfile *string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(*pkfile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, os.ErrInvalid
	}
	pkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pkey.(*rsa.PrivateKey), nil
}

func (c *CmdAppToken) GetInput() interface{} {
	return nil
}

func (c *CmdAppToken) ExecuteWithArgs(i interface{}) error {
	return nil
}

package auth

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"cli/cmd"
	"cli/common"
	"cli/config"
)

const (
	CmdAuthCodeName = "auth_code"
	AuthAppName     = "krypton_auth_code"
	ocPath          = "onecloud/token?app_name=krypton_auth_code&token_type=org-aware"
	hpbpPath        = "hpbp/token?feature=enroll"
)

type GetAuthTokenFlags struct {
	Code      *string
	State     *string
	TokenType *string
	URL       *string
}

type CmdAuthCode struct {
	cmd.CmdBase
	AuthBase
	AuthCodeInput
	BasePath  string
	TokenPath string
}

type AuthCodeInput struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

var (
	getAuthTokenFlags = NewGetAuthTokenFlags()
)

func init() {
	commands[CmdAuthCodeName] = NewCmdAuthCode()
}

func NewCmdAuthCode() *CmdAuthCode {
	c := CmdAuthCode{
		cmd.CmdBase{
			Name: CmdAuthCodeName,
		},
		AuthBase{},
		AuthCodeInput{},
		"services/oauth_handler",
		"",
	}
	fs := c.BaseInitFlags()
	(&c.AuthBase).initFlags(fs, AuthTokenTypeAuth)
	(&getAuthTokenFlags).initFlags(fs)
	return &c
}

func (f *GetAuthTokenFlags) initFlags(fs *flag.FlagSet) {
	fs.StringVar(f.URL, "url", "",
		"User can pass complete redirect url received from auth code flow to generate token.")
	// read code and state from flags if they are passed separately.
	fs.StringVar(f.Code, "code", "", "auth code")
	fs.StringVar(f.State, "state", "1111111111", "state")
}

func (f *GetAuthTokenFlags) verify() bool {
	if *f.Code == "" || *f.State == "" {
		return false
	}
	return true
}

func (f *GetAuthTokenFlags) parseURL() {
	// if urlArg is passed, extract auth code & state and use those. This will have first preference.
	if len(*f.URL) > 0 {
		u, err := url.Parse(*f.URL)
		if err != nil {
			log.Printf("error while parsing auth code redirect url. err - %v\n", err)
			return
		}
		*f.Code = u.Query().Get("code")
		*f.State = u.Query().Get("state")
	}
}

func (c *CmdAuthCode) Parse(args []string) (cmd.Command, error) {
	var err error
	c.BaseParse(args)
	getAuthTokenFlags.parseURL() // parse url if user passed.
	if !getAuthTokenFlags.verify() {
		log.Error("Missing required arguments [-code|-state]")
		log.Info("\n>>>>>>>>>>\nPlease login at below url to generate new auth code & state.\n" +
			getAuthCodeEntryPointUrl(*c.Flags.Server, *c.Flags.TokenType) + "\n" +
			"<<<<<<<<<<\n")
		return nil, cmd.ErrMissingArgs
	} else {
		c.Code = *getAuthTokenFlags.Code
		c.State = *getAuthTokenFlags.State
		if *getAuthTokenFlags.TokenType != "" {
			c.Flags.TokenType = getAuthTokenFlags.TokenType
		}
	}
	c.RunFunc = c.getToken

	if c.Stdin {
		err = cmd.ErrParseStdin
	}
	return c, err
}

// fetch token from oauth handler.
func (c *CmdAuthCode) getToken() {
	payload, err := json.Marshal(c.AuthCodeInput)
	if err != nil {
		log.Fatalf("cannot marshal request: %v\n", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.getTokenUrl(), bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("request error: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Content-Type", "application/json")

	log.HttpRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("enroll token request failed: %v\n", err)
	}
	log.HttpResponse(resp, nil)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("enroll token request failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read enroll token response. Error: %v\n", err)
	}
	log.HttpResponse(resp, result)
	// update cli cache with token
	if err = config.GetEnrollTokenCache().Update(result); err != nil {
		log.Debug("failed to cache access token")
	}
	// print token to console for user
	fmt.Println(string(result))
}

func (c *CmdAuthCode) GetInput() interface{} {
	return nil
}

func (c *CmdAuthCode) ExecuteWithArgs(i interface{}) error {
	return nil
}

func (c *CmdAuthCode) getTokenUrl() string {
	c.TokenPath = hpbpPath
	if *c.Flags.TokenType == TokenProviderOneCloud {
		c.TokenPath = ocPath
	}

	return fmt.Sprintf("%s/%s/%s", *c.Flags.Server, c.BasePath, c.TokenPath)
}

func NewGetAuthTokenFlags() GetAuthTokenFlags {
	return GetAuthTokenFlags{
		Code:      new(string),
		State:     new(string),
		TokenType: new(string),
		URL:       new(string),
	}
}

func getAuthCodeEntryPointUrl(server, tokenType string) string {
	url := fmt.Sprintf("%s/%s", server, "device-login/signin?auth_type=auth_code&state=1111111111")
	if tokenType == TokenProviderOneCloud {
		url = fmt.Sprintf("%s/%s", server, "device-login/signin?onecloud=true")
	}

	return url
}

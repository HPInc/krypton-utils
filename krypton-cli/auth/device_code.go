package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cli/cmd"
	"cli/common"
	"cli/config"
)

const (
	CmdDeviceCodeName = "device_code"
	TokenType         = "hpbp"
	deviceCodeAppName = "krypton_device_code"
)

type DeviceCodeAuth struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri_complete"`
}

type CmdDeviceCode struct {
	cmd.CmdBase
	AuthBase
	DeviceCodePath    string
	AuthorizationPath string
	TokenPath         string
}

func init() {
	commands[CmdDeviceCodeName] = NewCmdDeviceCode()
}

func NewCmdDeviceCode() *CmdDeviceCode {
	c := CmdDeviceCode{
		cmd.CmdBase{
			Name: CmdDeviceCodeName,
		},
		AuthBase{},
		"services/oauth_handler/device",
		"authorization",
		"token",
	}
	fs := c.BaseInitFlags()
	(&c.AuthBase).initFlags(fs, AuthTokenTypeAuth)
	return &c
}

func (c *CmdDeviceCode) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	c.RunFunc = c.deviceCodeFlow
	return c, nil
}

func (c *CmdDeviceCode) deviceCodeFlow() {
	data := url.Values{}
	data.Set("scope", "openid+enroll")

	req, err := http.NewRequest(http.MethodPost, c.getAuthorizationUrl(),
		strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("request error: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	log.HttpRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Device code request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Device code request failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read device code response. Error: %v\n", err)
	}
	log.HttpResponse(resp, result)

	dca := DeviceCodeAuth{}
	err = json.Unmarshal(result, &dca)
	if err != nil {
		log.Fatalf("Failed to unmarshal result: %v", err)
	}

	log.Println("Please login at this url: ", dca.VerificationUri)

	status := 0
	var i uint
	for i = 0; i < c.RetryCount; i++ {
		status = c.getToken(&dca)
		if status == http.StatusOK || status >= http.StatusInternalServerError {
			break
		} else {
			log.Debugf("Waiting for token acquire: %d / %d\n", i, c.RetryCount)
			time.Sleep(time.Second * time.Duration((*c.Flags.RetryDelay)))
		}
	}
	if status != http.StatusOK {
		log.Fatal("unable to acquire token")
	}
}

// wait till token is available or retries are exhausted
func (c *CmdDeviceCode) getToken(dca *DeviceCodeAuth) int {
	data := url.Values{}
	data.Set("grant_type", "uurn:ietf:params:oauth:grant-type:device_code")
	data.Set("device_code", dca.DeviceCode)

	req, err := http.NewRequest(http.MethodPost, c.getTokenUrl(),
		strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("request error: %v\n", err)
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.HttpRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("enroll token request failed: %v\n", err)
	}
	log.HttpResponse(resp, nil)
	// response is 400 for the first time if called fast enough
	if resp.StatusCode == 400 || resp.StatusCode == 429 {
		return resp.StatusCode
	}

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
	return resp.StatusCode
}

func (c *CmdDeviceCode) GetInput() interface{} {
	return nil
}

func (c *CmdDeviceCode) ExecuteWithArgs(i interface{}) error {
	return nil
}

func (c *CmdDeviceCode) getTokenUrl() string {
	path := c.TokenPath
	if *c.Flags.TokenType == TokenProviderOneCloud {
		path = fmt.Sprintf("%s?app_name=%s", path, deviceCodeAppName)
	}

	return fmt.Sprintf("%s/%s", c.getDeviceCodeUrl(), path)
}

func (c *CmdDeviceCode) getAuthorizationUrl() string {
	authPath := c.AuthorizationPath
	if *c.Flags.TokenType == TokenProviderOneCloud {
		authPath = fmt.Sprintf("%s?app_name=%s", authPath, deviceCodeAppName)
	}

	return fmt.Sprintf("%s/%s", c.getDeviceCodeUrl(), authPath)
}

func (c *CmdDeviceCode) getDeviceCodeUrl() string {
	return fmt.Sprintf("%s/%s", *c.Flags.Server, c.DeviceCodePath)
}

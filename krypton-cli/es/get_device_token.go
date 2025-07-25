package es

import (
	"errors"
	"fmt"

	"cli/cmd"
	"cli/common"
	"cli/config"
	"cli/es/enroll"
)

const (
	CMD_GET_DEVICE_TOKEN = "get_device_token"
)

var (
	ErrInvalidCertResponse = errors.New("returned certificate response is not valid")
)

type CmdGetDeviceToken struct {
	cmd.CmdBase
	EnrollBase
	CertDetails
}

type CmdGetDeviceTokenResult struct {
	DeviceToken string `json:"device_token"`
	DeviceId    string `json:"device_id"`
}

func init() {
	commands[CMD_GET_DEVICE_TOKEN] = NewCmdGetDeviceToken()
}

func NewCmdGetDeviceToken() *CmdGetDeviceToken {
	c := CmdGetDeviceToken{
		cmd.CmdBase{Name: CMD_GET_DEVICE_TOKEN},
		EnrollBase{},
		CertDetails{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&enrollFlags).initFlags(fs)
	return &c
}

func (c *CmdGetDeviceToken) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.Client.HardwareHash = *enrollFlags.hardwareHash
	c.Client.BulkEnrollToken = *enrollFlags.bulkEnrollToken
	c.RunFunc = c.getDeviceToken
	var err error
	if c.Stdin {
		err = cmd.ErrParseStdin
	}
	return c, err
}

// main worker function of the command
func (c *CmdGetDeviceToken) getDeviceToken() {
	// make a copy of the client first
	client := c.Client

	var cr *enroll.CertificateResponse
	var err error

	if c.Stdin {
		cr = &enroll.CertificateResponse{
			DeviceId:    c.DeviceId,
			Certificate: c.Certificate,
		}
		client.PK, err = common.PrivateKeyFromPem(c.PrivateKey)
		if err != nil {
			log.Fatal("Error: ", err)
		}
	} else {
		cr, err = client.GetDeviceCertificate()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
	log.Debug("Getting device token for:", cr.DeviceId)
	dt, err := client.GetDeviceToken(cr)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	r := CmdGetDeviceTokenResult{
		DeviceToken: dt,
		DeviceId:    cr.DeviceId,
	}

	result := common.GetJsonString(r)

	// update cli cache with device_token
	// only write cache on single count iterations
	if c.Count == 1 {
		if err = config.GetDeviceTokenCache().Update([]byte(result)); err != nil {
			log.Debug("failed to cache device token")
		}
	}
	fmt.Println(result)
}

func (c *CmdGetDeviceToken) GetInput() interface{} {
	return &c.CertDetails
}

func (c *CmdGetDeviceToken) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}

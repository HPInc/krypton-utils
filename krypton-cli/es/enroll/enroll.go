package enroll

import (
	"bytes"
	"cli/common"
	"cli/logging"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall"
	"time"
)

type TokenType string

const (
	TokenTypeAzureAD    TokenType = "azuread"
	TokenTypeDevice     TokenType = "device"
	TokenTypeEnrollment TokenType = "enrollment"
	TokenTypeHpbp       TokenType = "hpbp"

	Bearer          = "Bearer"
	HeaderTokenType = "X-HP-Token-Type" //#nosec G101
	RetryAfter      = "Retry-After"

	DefaultRetryCount        = 3
	DefaultRetryAfterSeconds = 5
	BearerTokenRetryCount    = 10
	EnrollRetryCount         = 10
	RetrySecondsMultiplier   = 2
)

type enrollRequest struct {
	CSR               string `json:"csr"`
	ManagementService string `json:"mgmt_service"`
	HardwareHash      string `json:"hardware_hash"`
}

type EnrollResponse struct {
	Id     string `json:"id"`
	Bearer string `json:"-"`
}

type CertificateResponse struct {
	Id            string `json:"id"`
	DeviceId      string `json:"device_id,omitempty"`
	Certificate   string `json:"certificate,omitempty"`
	ReceiptHandle string `json:"receipt_handle,omitempty"`
	IsValid       bool   `json:"-"`
	RetryAfter    int    `json:"retry_after,omitempty"`
	Status        string `json:"status,omitempty"`
}

type EnrollClient struct {
	EnrollUrl        string          `json:"server_url"`
	DSTSServer       string          `json:"dsts_server"`
	RetryCount       uint            `json:"retry_count"`
	TokenServer      string          `json:"token_server"`
	TokenType        string          `json:"token_type"`
	PK               *rsa.PrivateKey `json:"-"`
	JWTToken         string          `json:"jwt_token"`
	HardwareHash     string          `json:"hardware_hash"`
	BulkEnrollToken  string          `json:"bulk_enroll_token"`
	ManagementServer string          `json:"mgmt_server"`
}

type EnrollmentTokenResponse struct {
	Token string `json:"token"`
}

var log = logging.GetLogger()

// set token type. primarily to allow clients to alter
// the token type on different usages like renew enroll.
func (c *EnrollClient) SetTokenType(tokenType string) {
	c.TokenType = tokenType
	if tokenType == "enrollment" {
		c.TokenServer = fmt.Sprintf("%s/enrollmenttoken", c.DSTSServer)
	}
}

func (c *EnrollClient) SetRetryCount(retryCount uint) {
	c.RetryCount = retryCount
	if c.RetryCount == 0 {
		c.RetryCount = DefaultRetryCount
	}
}

func (c *EnrollClient) GetDeviceCertificate() (*CertificateResponse, error) {
	var cr *CertificateResponse
	er, err := c.Enroll()
	if err != nil {
		return nil, err
	}
	cr, err = c.getEnrollStatus(er)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (c *EnrollClient) Enroll() (*EnrollResponse, error) {
	bearer, err := c.getBearerTokenWithRetry("")
	if err != nil {
		log.Printf("Failed to get bearer token. Error: %v\n", err)
		return nil, err
	}
	jsonString, err := c.getEnrollPayload()
	if err != nil {
		log.Printf("Failed to initialize enrollment payload. Error: %v\n", err)
		return nil, err
	}
	enrollUrl := fmt.Sprintf("%s/enroll", c.EnrollUrl)
	log.Debug("enroll:", c.EnrollUrl, c.TokenType, string(jsonString))
	req, err := http.NewRequest("POST", enrollUrl, bytes.NewBuffer(jsonString))
	if err != nil {
		log.Printf("Failed to create enrollment request. Error: %v\n", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add(HeaderTokenType, c.TokenType)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	return getEnrollResponse(&bearer, req)
}

// retries and gets enroll response
func getEnrollResponse(bearer *string, req *http.Request) (*EnrollResponse, error) {
	client := common.RetriableClient(EnrollRetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send enrollment request. Retries: %d, Error: %v\n",
			EnrollRetryCount, err)
		return nil, err
	}
	log.HttpResponse(resp, nil)
	if resp.StatusCode != http.StatusAccepted {
		log.Printf("Enrollment request failed with error code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("enrollment request failed")
	}

	er := EnrollResponse{Bearer: *bearer}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read enrollment response. Error: %v\n", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	err = json.Unmarshal(data, &er)
	if err != nil {
		log.Printf("Failed to unmarshal enrollment response. Error: %v\n", err)
		return nil, err
	}
	return &er, nil
}

func (c *EnrollClient) getBearerTokenWithRetry(tenantId string) (string, error) {
	var token string
	var err error
	for i := 1; i <= BearerTokenRetryCount; i++ {
		token, err = c.getBearerToken(tenantId)
		if err != nil {
			// only retry for connrefused
			if !errors.Is(err, syscall.ECONNREFUSED) {
				break
			}
			time.Sleep(time.Second * time.Duration(i))
			continue
		}
		return token, nil
	}
	return "", err
}

func (c *EnrollClient) getBearerToken(tenantId string) (string, error) {
	switch TokenType(c.TokenType) {
	case TokenTypeAzureAD:
		return c.getAzureADToken(tenantId)
	case TokenTypeEnrollment:
		return c.getEnrollmentToken(tenantId)
	case TokenTypeHpbp:
		return c.getJwtToken(), nil
	case TokenTypeDevice:
		return c.getJwtToken(), nil
	default:
		return "", fmt.Errorf("invalid token type: %s", c.TokenType)
	}
}

// fetch enrollment token by reaching out to configured server
func (c *EnrollClient) fetchEnrollmentToken(tenantId string) (string, error) {
	url := fmt.Sprintf("%s?tenant_id=%s", c.TokenServer, tenantId)
	resp, err := common.NewClient().Post(url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	tokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read bearer token from response. Error: %v\n", err)
		return "", err
	}
	et := EnrollmentTokenResponse{}
	if err = json.Unmarshal(tokenBytes, &et); err != nil {
		log.Printf("Failed to unmarshal token. Error: %v\n", err)
		return "", err
	}
	token := "Bearer " + et.Token
	return strings.TrimSpace(token), nil
}

func (c *EnrollClient) getAzureADToken(tenantId string) (string, error) {
	if c.JWTToken != "" {
		return "Bearer " + c.JWTToken, nil
	}
	url := c.TokenServer
	if tenantId != "" {
		url = fmt.Sprintf("%s?tenant_id=%s", url, tenantId)
	}
	resp, err := common.NewClient().Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	tokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read bearer token from response. Error: %v\n", err)
		return "", err
	}
	token := "Bearer " + string(tokenBytes)
	return strings.TrimSpace(token), nil
}

// device code flow, not self contained, requires external approval
// this flow will require HTTP_PROXY now but it can be externally
// controlled via the same env var so we do not provide extra support here
func (c *EnrollClient) getJwtToken() string {
	return fmt.Sprintf("%s %s", Bearer, c.JWTToken)
}

// enrollment path. use JWTToken if present
func (c *EnrollClient) getEnrollmentToken(tenantId string) (string, error) {
	if c.BulkEnrollToken == "" {
		return c.fetchEnrollmentToken(tenantId)
	} else {
		return fmt.Sprintf("%s %s", Bearer, c.BulkEnrollToken), nil
	}
}

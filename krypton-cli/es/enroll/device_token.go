package enroll

import (
	"cli/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	PARAM_CLIENT_ASSERTION      = "client_assertion"
	PARAM_CLIENT_ASSERTION_TYPE = "client_assertion_type"
	CLIENT_ASSERTION_TYPE       = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
)

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

type DeviceTokenResponse struct {
	Token string `json:"access_token"`
}

type AssertionClaims struct {
	// Standard JWT claims such as 'aud', 'exp', 'jti', 'iat', 'iss', 'nbf',
	// 'sub'
	jwt.RegisteredClaims

	// Nonce - the challenge returned by the device STS challenge api
	Nonce string `json:"nonce"`
}

// flow to get a device token
// device id is obtained from enroll
// see https://rndwiki.inc.hpicorp.net/confluence/display/Hptm01/DSTS+API%3A+Device+authentication
func (c *EnrollClient) GetDeviceToken(cr *CertificateResponse) (string, error) {
	// get challenge code
	challenge, err := c.getChallengeCode(cr.DeviceId)
	if err != nil {
		log.Printf("DeviceToken: Failed to get challenge code. DeviceId: %s, Error: %v\n",
			cr.DeviceId, err)
		return "", err
	}
	assertion, err := c.getAssertion(challenge, cr)
	if err != nil {
		log.Printf("Failed to generate signed client assertion. DeviceId: %s, Error: %v",
			cr.DeviceId, err)
		return "", nil
	}
	log.Debug("assertion: ", assertion)
	return c.getDeviceTokenFromSTS(assertion)
}

func (c *EnrollClient) getAssertion(challenge string, cr *CertificateResponse) (string, error) {
	claims := getAssertionClaims(challenge, cr.DeviceId)
	assertionToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	assertionToken.Header["x5c"] = []string{cr.Certificate}
	assertion, err := assertionToken.SignedString(c.PK)
	if err != nil {
		log.Printf("Failed to generate signed client assertion. DeviceId: %s, Error: %v",
			cr.DeviceId, err)
		return "", err
	}
	return assertion, nil
}

// get device token from sts
func (c *EnrollClient) getDeviceTokenFromSTS(assertion string) (string, error) {
	data := url.Values{}
	data.Set(PARAM_CLIENT_ASSERTION_TYPE, CLIENT_ASSERTION_TYPE)
	data.Set(PARAM_CLIENT_ASSERTION, assertion)

	tokenUrl := fmt.Sprintf("%s/deviceauth/token", c.DSTSServer)
	req, err := http.NewRequest(http.MethodPost, tokenUrl,
		strings.NewReader(data.Encode()))
	if err != nil {
		log.Errorf("Failed to create device token request. Error: %v\n",
			err)
		return "", err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute device token request. Error: %v\n", err)
		return "", err
	}
	return getDeviceTokenFromResponse(resp)
}

func getDeviceTokenFromResponse(r *http.Response) (string, error) {
	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"Device token request error: %d", r.StatusCode)
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read token response. Error: %v\n", err)
		return "", err
	}
	log.HttpResponse(r, data)
	dtr := DeviceTokenResponse{}
	err = json.Unmarshal(data, &dtr)
	if err != nil {
		log.Printf("Failed to unmarshal token response. Error: %v\n", err)
		return "", err
	}
	return dtr.Token, nil
}

func (c *EnrollClient) getChallengeCode(deviceId string) (string, error) {
	queryUrl := fmt.Sprintf("%s/deviceauth/challenge?device_id=%s",
		c.DSTSServer, deviceId)
	log.Debug("getChallengeCode: ", queryUrl)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		log.Printf("Failed to create challenge request. DeviceId: %s, Error: %v\n",
			deviceId, err)
		return "", err
	}
	common.AddUserAgentHeader(req)
	client := common.RetriableClient(DefaultRetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute challenge request. DeviceId: %s, Error: %v\n",
			deviceId, err)
		return "", err
	}
	return getChallengeCodeResponse(resp)
}

func getChallengeCodeResponse(r *http.Response) (string, error) {
	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"Challenge code request returned: %d", r.StatusCode)
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read challenge response. Error: %v\n", err)
		return "", err
	}
	log.HttpResponse(r, data)
	cr := ChallengeResponse{}
	err = json.Unmarshal(data, &cr)
	log.Debug(cr)
	if err != nil {
		log.Printf("Failed to unmarshal challenge response. Error: %v\n", err)
		return "", err
	}
	return cr.Challenge, nil
}

// Construct assertion claims struct
func getAssertionClaims(challenge, deviceId string) *AssertionClaims {
	return &AssertionClaims{
		Nonce: challenge,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    deviceId,
			Subject:   deviceId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
}

package auth

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
	paramAppId               = "app_id"
	paramClientAssertion     = "client_assertion"
	paramClientAssertionType = "client_assertion_type"
	clientAssertionType      = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
)

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

type AppTokenResponse struct {
	Token string `json:"access_token"`
}

type AssertionClaims struct {
	// Standard JWT claims such as 'aud', 'exp', 'jti', 'iat', 'iss', 'nbf',
	// 'sub'
	jwt.RegisteredClaims

	// Nonce - the challenge returned by the device STS challenge api
	Nonce string `json:"nonce"`
}

func (c *CmdAppToken) getAssertion(challenge string) (string, error) {
	claims := getAssertionClaims(challenge, *c.AppId)
	assertionToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	assertion, err := assertionToken.SignedString(c.PrivateKey)
	if err != nil {
		log.Printf("Failed to generate signed client assertion. AppId: %s, Error: %v",
			*c.AppId, err)
		return "", err
	}
	return assertion, nil
}

// get app token from dsts
func (c *CmdAppToken) getAppTokenFromSTS(assertion string) (string, error) {
	data := url.Values{}
	data.Set(paramClientAssertionType, clientAssertionType)
	data.Set(paramClientAssertion, assertion)
	data.Set(paramAppId, *c.AppId)

	tokenUrl := fmt.Sprintf("%s/%s/appauth/token", *c.Flags.Server, c.ApiBasePath)
	req, _ := http.NewRequest(http.MethodPost, tokenUrl,
		strings.NewReader(data.Encode()))
	common.AddUserAgentHeader(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute app token request. Error: %v\n", err)
		return "", err
	}
	return getAppTokenFromResponse(resp)
}

func getAppTokenFromResponse(r *http.Response) (string, error) {
	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"App token request error: %d", r.StatusCode)
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read token response. Error: %v\n", err)
		return "", err
	}
	log.HttpResponse(r, data)
	aptr := AppTokenResponse{}
	err = json.Unmarshal(data, &aptr)
	if err != nil {
		log.Printf("Failed to unmarshal token response. Error: %v\n", err)
		return "", err
	}
	return aptr.Token, nil
}

func (c *CmdAppToken) getChallengeCode() (string, error) {
	queryUrl := fmt.Sprintf("%s/%s/appauth/challenge?%s=%s",
		*c.Flags.Server, c.ApiBasePath, paramAppId, *c.AppId)
	log.Debug("getChallengeCode: ", queryUrl)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		log.Printf("Failed to create challenge request. AppId: %s, Error: %v\n",
			*c.AppId, err)
		return "", err
	}
	common.AddUserAgentHeader(req)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute challenge request. AppId: %s, Error: %v\n",
			*c.AppId, err)
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
	if err != nil {
		log.Printf("Failed to unmarshal challenge response. Error: %v\n", err)
		return "", err
	}
	return cr.Challenge, nil
}

// Construct assertion claims struct
func getAssertionClaims(challenge, appId string) *AssertionClaims {
	return &AssertionClaims{
		Nonce: challenge,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    appId,
			Subject:   appId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
}

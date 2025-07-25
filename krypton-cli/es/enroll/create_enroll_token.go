package enroll

import (
	"cli/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EnrollTokenResponse struct {
	TenantId  string `json:"tenant_id,omitempty"`
	Token     string `json:"enroll_token,omitempty"`
	IssuedAt  int64  `json:"issued_at,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	HttpCode  int    `json:"http_code"`
}

func (c *EnrollClient) CreateEnrollToken() (*EnrollTokenResponse, error) {
	bearer, err := c.getBearerTokenWithRetry("")
	if err != nil {
		log.Printf("Failed to get bearer token. Error: %v\n", err)
		return nil, err
	}
	queryUrl := fmt.Sprintf("%s/enroll_token", c.EnrollUrl)
	log.Debug("get enroll token:", queryUrl)
	req, err := http.NewRequest("POST", queryUrl, nil)
	if err != nil {
		log.Printf("Failed to create enroll token request. Error: %v\n", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add(HeaderTokenType, c.TokenType)
	req.Header.Add("Authorization", bearer)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Get enroll token failed: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.HttpResponse(resp, nil)
		log.Fatalf("Create enroll token failed with error code: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read enrollment status check response. Error: %v\n", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	er := EnrollTokenResponse{HttpCode: resp.StatusCode}
	err = json.Unmarshal(data, &er)
	if err != nil {
		log.Printf("Failed to unmarshal enroll token response. Error: %v\n", err)
		return nil, err
	}
	return &er, nil
}

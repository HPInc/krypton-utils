package enroll

import (
	"cli/common"
	"fmt"
	"net/http"
)

func (c *EnrollClient) DeleteEnrollToken() (*EnrollTokenResponse, error) {
	bearer, err := c.getBearerTokenWithRetry("")
	if err != nil {
		log.Printf("Failed to get bearer token. Error: %v\n", err)
		return nil, err
	}
	queryUrl := fmt.Sprintf("%s/enroll_token", c.EnrollUrl)
	log.Debug("get enroll token:", queryUrl)
	req, err := http.NewRequest("DELETE", queryUrl, nil)
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
		log.Fatalf("Delete enroll token failed: %v\n", err)
	}
	er := EnrollTokenResponse{HttpCode: resp.StatusCode}
	if resp.StatusCode != http.StatusOK {
		log.Error("Delete enroll token failed with error code: ", resp.StatusCode)
	}
	log.HttpResponse(resp, nil)
	return &er, nil
}

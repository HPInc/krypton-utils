package enroll

import (
	"cli/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *EnrollClient) GetStatus(id string) (*CertificateResponse, error) {
	bearer, err := c.getBearerTokenWithRetry("")
	log.Info(bearer, err)
	if err != nil {
		log.Printf("Failed to get bearer token. Error: %v\n", err)
		return nil, err
	}
	return c.getEnrollStatus(&EnrollResponse{Bearer: bearer, Id: id})
}

func (c *EnrollClient) getEnrollStatus(er *EnrollResponse) (*CertificateResponse, error) {
	queryUrl := fmt.Sprintf("%s/enroll/%s", c.EnrollUrl, er.Id)
	log.Debug("get enroll status:", c.EnrollUrl, er.Id)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		log.Printf("Failed to create enrollment status check request. Error: %v\n", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add(HeaderTokenType, c.TokenType)
	req.Header.Add("Authorization", er.Bearer)
	client := common.RetriableClient(c.RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Get enrollment status check request. Error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read enrollment status check response. Error: %v\n", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	cr := CertificateResponse{IsValid: true}
	err = json.Unmarshal(data, &cr)
	if err != nil {
		log.Printf("Failed to unmarshal enrollment status check response. Error: %v\n", err)
		return nil, err
	}
	return &cr, nil
}

package enroll

import (
	"bytes"
	"cli/common"
	"fmt"
	"net/http"

	uuid "github.com/google/uuid"
)

func (c *EnrollClient) RenewEnroll(
	deviceId uuid.UUID,
	deviceToken string) (*EnrollResponse, error) {
	bearer := fmt.Sprintf("Bearer %s", deviceToken)
	jsonString, err := c.getEnrollPayload()
	if err != nil {
		fmt.Printf("Failed to initialize enrollment payload. Error: %v\n", err)
		return nil, err
	}
	enrollUrl := fmt.Sprintf("%s/enroll/%s", c.EnrollUrl, deviceId)
	req, err := http.NewRequest(http.MethodPatch, enrollUrl, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Printf("Failed to renew enrollment request. Error: %v\n", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add(HeaderTokenType, string(TokenTypeDevice))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	return getEnrollResponse(&bearer, req)
}

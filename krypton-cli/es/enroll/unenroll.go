package enroll

import (
	"bytes"
	"cli/common"
	"fmt"
	"net/http"

	uuid "github.com/google/uuid"
)

func (c *EnrollClient) Unenroll(
	deviceId uuid.UUID,
	deviceToken string) (*EnrollResponse, error) {
	bearer := fmt.Sprintf("Bearer %s", deviceToken)
	jsonString, err := c.getEnrollPayload()
	if err != nil {
		fmt.Printf("Failed to initialize unenrollment payload. Error: %v\n", err)
		return nil, err
	}
	enrollUrl := fmt.Sprintf("%s/enroll/%s", c.EnrollUrl, deviceId)
	req, err := http.NewRequest(http.MethodDelete, enrollUrl, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Printf("Failed to unenrollment request. Error: %v\n", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add(HeaderTokenType, string(TokenTypeDevice))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	return getEnrollResponse(&bearer, req)
}

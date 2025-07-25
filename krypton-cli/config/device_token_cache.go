package config

import (
	"cli/common"
	"encoding/json"
	"fmt"
)

const (
	DeviceTokenFile = "device_token"
	ClaimTenantId   = "tid"
)

var (
	// cache file path
	DeviceCacheFile = fmt.Sprintf("%s/%s", getConfigDir(), DeviceTokenFile)
)

// cached token
type deviceTokenData struct {
	DeviceToken string `json:"device_token"`
	DeviceId    string `json:"device_id"`
}

type DeviceTokenCache struct {
	tokenData *deviceTokenData
}

func NewDeviceTokenCache() *DeviceTokenCache {
	return &DeviceTokenCache{
		tokenData: readCachedDeviceData(),
	}
}

// set device token
func (e *DeviceTokenCache) Update(bytes []byte) error {
	if err := common.WriteFile(DeviceCacheFile, bytes); err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, e.tokenData); err != nil {
		return err
	}
	return nil
}

// read cached data if any
func readCachedDeviceData() *deviceTokenData {
	et := deviceTokenData{}
	bytes, err := common.ReadFile(DeviceCacheFile)
	if err != nil {
		logger.Debugf("Error reading device token from cache: file=%s, err=%v",
			DeviceCacheFile, err)
		return nil
	}
	if err = json.Unmarshal(bytes, &et); err != nil {
		logger.Debugf("Error unmarshal token from cache: file=%s, err=%v",
			DeviceCacheFile, err)
		return nil
	}
	return &et
}

// get token data
func (e *DeviceTokenCache) GetDeviceToken() string {
	if e.tokenData != nil {
		return e.tokenData.DeviceToken
	}
	return ""
}

// get device id
func (e *DeviceTokenCache) GetDeviceId() string {
	if e.tokenData != nil {
		return e.tokenData.DeviceId
	}
	return ""
}

// get tid claim from token
func (e *DeviceTokenCache) GetTenantId() string {
	if e.tokenData == nil || e.tokenData.DeviceToken == "" {
		return ""
	}
	val, err := common.GetClaimFromToken(e.tokenData.DeviceToken, ClaimTenantId)
	if err != nil {
		logger.Fatalf("no claim %s in token", ClaimTenantId)
	}
	return val
}

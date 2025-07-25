package config

import (
	"cli/common"
	"encoding/json"
	"fmt"
)

const (
	EnrollTokenFile = "enroll_token"
)

var (
	// cache file path
	EnrollCacheFile = fmt.Sprintf("%s/%s", getConfigDir(), EnrollTokenFile)
)

// cached token
type enrollTokenData struct {
	AccessToken string `json:"access_token"`
}

type EnrollTokenCache struct {
	tokenData *enrollTokenData
}

func NewEnrollTokenCache() *EnrollTokenCache {
	return &EnrollTokenCache{
		tokenData: readCachedEnrollData(),
	}
}

// set enroll token
func (e *EnrollTokenCache) Update(bytes []byte) error {
	if err := common.WriteFile(EnrollCacheFile, bytes); err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, e.tokenData); err != nil {
		return err
	}
	return nil
}

// read cached data if any
func readCachedEnrollData() *enrollTokenData {
	et := enrollTokenData{}
	bytes, err := common.ReadFile(EnrollCacheFile)
	if err != nil {
		logger.Debugf("Error reading enroll token from cache: file=%s, err=%v",
			EnrollCacheFile, err)
		return nil
	}
	if err = json.Unmarshal(bytes, &et); err != nil {
		logger.Debugf("Error unmarshal token from cache: file=%s, err=%v",
			EnrollCacheFile, err)
		return nil
	}
	return &et
}

// get token data
func (e *EnrollTokenCache) GetAccessToken() string {
	if e.tokenData != nil {
		return e.tokenData.AccessToken
	}
	return ""
}
